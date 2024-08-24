package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pawpawchat/s3/api/pb"
)

func (s *s3Server) PresignPutObject(ctx context.Context, req *pb.PresignPutObjectRequest) (*pb.PresignPutObjectResponse, error) {
	objectInfo, err := parsePresignPutObjectRequest(req)
	if err != nil {
		return nil, err
	}

	sign, err := s.presign.PresignPutObject(ctx, objectInfo)
	if err != nil {
		return nil, err
	}

	return &pb.PresignPutObjectResponse{URL: sign.URL}, nil
}

func parsePresignPutObjectRequest(req *pb.PresignPutObjectRequest) (*s3.PutObjectInput, error) {
	if req.OwnerId == 0 {
		return nil, fmt.Errorf("missing owner id")
	}

	if req.ObjKey == "" {
		return nil, fmt.Errorf("missing object key")
	}

	var expectedType pb.ContentType
	var ext string

	switch fileExt := req.GetFileExt().(type) {
	case *pb.PresignPutObjectRequest_AudioExt:
		ext = fileExt.AudioExt.String()
		expectedType = pb.ContentType_Audio

	case *pb.PresignPutObjectRequest_VideoExt:
		ext = fileExt.VideoExt.String()
		expectedType = pb.ContentType_Video

	case *pb.PresignPutObjectRequest_ImageExt:
		ext = fileExt.ImageExt.String()
		expectedType = pb.ContentType_Image

	default:
		return nil, fmt.Errorf("unknown file extension")
	}

	if req.ContentType != expectedType {
		return nil, fmt.Errorf("%s isn't the correct extension for content type %v", ext, req.ContentType)
	}

	key := fmt.Sprintf("%d/%v/%s.%s", req.OwnerId, req.ContentType, req.ObjKey, ext)
	bucket := "pawpawchat"
	mimeType := fmt.Sprintf("%s/%s", strings.ToLower(req.ContentType.String()), ext)

	return &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		ContentType: aws.String(mimeType),
	}, nil
}
