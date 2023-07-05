package grpcservice

import (
	"context"

	"github.com/ekkinox/fx-template/internal/model"
	"github.com/ekkinox/fx-template/internal/repository"
	"github.com/ekkinox/fx-template/modules/fxgrpcserver"
	"github.com/ekkinox/fx-template/modules/fxlogger"
	"github.com/ekkinox/fx-template/proto/posts"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type PostsCrudServer struct {
	posts.UnimplementedPostCrudServiceServer
	repository *repository.PostRepository
}

func NewPostsCrudServer(repository *repository.PostRepository, logger *fxlogger.Logger) *PostsCrudServer {
	return &PostsCrudServer{
		repository: repository,
	}
}

func (s *PostsCrudServer) GetPost(ctx context.Context, in *posts.GetPostRequest) (*posts.GetPostResponse, error) {

	fxgrpcserver.CtxLogger(ctx).Info().Msg("got grpc GetPost request")

	dbPost, err := s.repository.Find(ctx, int(in.Id.Value))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "post with id %d not found: %v", in.Id.Value, err)
	}

	return &posts.GetPostResponse{
		Success: true,
		Post:    modelToProto(dbPost),
	}, nil
}

func (s *PostsCrudServer) CreatePost(ctx context.Context, in *posts.CreatePostRequest) (*posts.CreatePostResponse, error) {

	fxgrpcserver.CtxLogger(ctx).Info().Msg("got grpc CreatePost request")

	dbPostCreate := new(model.Post)
	dbPostCreate.Title = in.Post.Title.GetValue()
	dbPostCreate.Description = in.Post.Description.GetValue()
	dbPostCreate.Likes = int(in.Post.Likes.GetValue())

	err := s.repository.Create(ctx, dbPostCreate)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while creating post: %v", err)
	}

	return &posts.CreatePostResponse{
		Success: true,
		Post:    modelToProto(dbPostCreate),
	}, nil
}

func (s *PostsCrudServer) UpdatePost(ctx context.Context, in *posts.UpdatePostRequest) (*posts.UpdatePostResponse, error) {

	fxgrpcserver.CtxLogger(ctx).Info().Msg("got grpc UpdatePost request")

	dbPost, err := s.repository.Find(ctx, int(in.Post.Id.Value))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "post with id %d not found: %v", in.Post.Id.Value, err)
	}

	dbPostUpdate := new(model.Post)
	if in.Post.Title.GetValue() != "" {
		dbPostUpdate.Title = in.Post.Title.GetValue()
	}
	if in.Post.Description.GetValue() != "" {
		dbPostUpdate.Description = in.Post.Description.GetValue()
	}
	if in.Post.Likes.GetValue() != 0 {
		dbPostUpdate.Likes = int(in.Post.Likes.GetValue())
	}

	err = s.repository.Update(ctx, dbPost, dbPostUpdate)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while updating post with id %d: %v", in.Post.Id.Value, err)
	}

	return &posts.UpdatePostResponse{
		Success: true,
		Post:    modelToProto(dbPost),
	}, nil
}

func (s *PostsCrudServer) DeletePost(ctx context.Context, in *posts.DeletePostRequest) (*posts.DeletePostResponse, error) {

	fxgrpcserver.CtxLogger(ctx).Info().Msg("got grpc DeletePost request")

	dbPost, err := s.repository.Find(ctx, int(in.Id.Value))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "post with id %d not found: %v", in.Id.Value, err)
	}

	err = s.repository.Delete(ctx, dbPost)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while deleting post with id %d: %v", in.Id.Value, err)
	}

	return &posts.DeletePostResponse{
		Success: true,
	}, nil
}

func (s *PostsCrudServer) ListPosts(ctx context.Context, in *emptypb.Empty) (*posts.ListPostsResponse, error) {

	fxgrpcserver.CtxLogger(ctx).Info().Msg("got grpc ListPosts request")

	dbPosts, err := s.repository.FindAll(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while listing posts: %v", err)
	}

	var postsList []*posts.Post

	for _, dbPost := range dbPosts {
		postsList = append(postsList, modelToProto(&dbPost))
	}

	return &posts.ListPostsResponse{
		Success: true,
		Posts: &posts.PostsList{
			Posts: postsList,
		},
	}, nil
}

func modelToProto(post *model.Post) *posts.Post {
	return &posts.Post{
		Id:          wrapperspb.Int32(int32(post.ID)),
		Title:       wrapperspb.String(post.Title),
		Description: wrapperspb.String(post.Description),
		Likes:       wrapperspb.Int32(int32(post.Likes)),
	}
}
