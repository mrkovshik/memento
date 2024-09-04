package grpcServer

import (
	"fmt"
	"io"
	"os"

	"github.com/mrkovshik/memento/internal/auth"
	"github.com/mrkovshik/memento/internal/model/data"
	pb "github.com/mrkovshik/memento/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) AddVariousData(stream pb.Memento_AddVariousDataServer) error {
	ctx := stream.Context()
	userID, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return stream.SendAndClose(&pb.AddVariousDataResponse{
			UploadStatus: &pb.UploadStatus{
				Success: false,
				Message: "Failed to get user ID",
			},
			Error: err.Error(),
		})
	}
	req, err := stream.Recv()
	if err != nil {
		return stream.SendAndClose(&pb.AddVariousDataResponse{
			UploadStatus: &pb.UploadStatus{
				Success: false,
				Message: "Failed to receive data",
			},
			Error: err.Error(),
		})
	}

	variousData := req.GetVariousData()
	if variousData == nil {
		return stream.SendAndClose(&pb.AddVariousDataResponse{
			UploadStatus: &pb.UploadStatus{
				Success: false,
				Message: "data model is required",
			},
			Error: "data model is empty",
		})
	}

	// Inserting entry to DB
	dataModel, err := s.service.AddVariousData(ctx, data.VariousData{
		UserID:   userID,
		DataType: int(variousData.DataType),
		Meta:     variousData.Meta,
	})
	if err != nil {
		return stream.SendAndClose(&pb.AddVariousDataResponse{
			UploadStatus: &pb.UploadStatus{
				Success: false,
				Message: "saving data failed",
			},
			Error: fmt.Sprintf("saving data failed: %s", err.Error()),
		})
	}

	// Prepare to receive file chunks
	dirName := fmt.Sprintf("data/%d", userID)
	fileName := fmt.Sprintf("./%s/%s", dirName, dataModel.UUID)
	if err := os.MkdirAll(dirName, os.ModePerm); err != nil {
		return stream.SendAndClose(&pb.AddVariousDataResponse{
			UploadStatus: &pb.UploadStatus{
				Success: false,
				Message: "failed to create or open file",
			},
			Error: err.Error(),
		})
	}
	dataFile, err := os.OpenFile(fileName, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return stream.SendAndClose(&pb.AddVariousDataResponse{
			UploadStatus: &pb.UploadStatus{
				Success: false,
				Message: "failed to create or open file",
			},
			Error: err.Error(),
		})
	}
	defer dataFile.Close()

	// Receiving file by chunks
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// File transmission is complete
			break
		}
		if err != nil {
			return stream.SendAndClose(&pb.AddVariousDataResponse{
				UploadStatus: &pb.UploadStatus{
					Success: false,
					Message: "upload file error",
				},
				Error: err.Error(),
			})
		}

		if chunk := req.GetChunk(); chunk != nil {
			if _, err := dataFile.Write(chunk); err != nil {
				return stream.SendAndClose(&pb.AddVariousDataResponse{
					UploadStatus: &pb.UploadStatus{
						Success: false,
						Message: "failed to write chunk to file",
					},
					Error: err.Error(),
				})
			}
		}
	}

	return stream.SendAndClose(&pb.AddVariousDataResponse{
		UploadStatus: &pb.UploadStatus{
			Success: true,
			Message: "Data saved successfully!",
		},
	})
}

func (s *Server) DownloadVariousDataFile(req *pb.DownloadVariousDataFileRequest, stream pb.Memento_DownloadVariousDataFileServer) error {
	ctx := stream.Context()
	userID, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "user ID is missing: %v", err)
	}
	// Construct the file path based on the file ID
	filePath := fmt.Sprintf("./data/%d/%s", userID, req.DataUUID)

	// Open the file for reading
	file, err := os.Open(filePath)
	if err != nil {
		return status.Errorf(codes.NotFound, "file not found: %v", err)
	}
	defer file.Close()

	// Define a buffer size for reading the file in chunks
	const chunkSize = 1024 * 1024 // 1MB chunks
	buffer := make([]byte, chunkSize)

	for {
		// Read a chunk from the file
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return status.Errorf(codes.Internal, "failed to read file: %v", err)
		}

		if n == 0 {
			break // Reached the end of the file
		}

		// Send the chunk to the client
		chunk := buffer[:n] // Only send the bytes read

		if err := stream.Send(&pb.DownloadVariousDataFileResponse{Chunk: chunk}); err != nil {
			return status.Errorf(codes.Internal, "failed to send chunk: %v", err)
		}

		if err == io.EOF {
			break
		}
	}

	return nil
}
