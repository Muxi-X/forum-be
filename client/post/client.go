package main

import (
	"context"
	"fmt"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/handler"
	tracer "forum/pkg/tracer"
	"github.com/micro/go-micro"
	"github.com/opentracing/opentracing-go"
	"log"

	opentracingWrapper "github.com/micro/go-plugins/wrapper/trace/opentracing"
)

func main() {
	t, io, err := tracer.NewTracer("forum.service.post", "localhost:6831")
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	defer logger.SyncLogger()

	opentracing.SetGlobalTracer(t)

	service := micro.NewService(micro.Name("forum.cli.post"),
		micro.WrapClient(
			opentracingWrapper.NewClientWrapper(opentracing.GlobalTracer()),
		),
		micro.WrapCall(handler.ClientErrorHandlerWrapper()),
	)

	service.Init()

	client := pb.NewPostServiceClient("forum.service.post", service.Client())
	//
	// _, err = client.CreatePost(context.TODO(), &pb.CreatePostRequest{
	// 	UserId:   2,
	// 	Content:  "外比巴卜",
	// 	TypeId:   1, // 默认为1
	// 	Title:    "first post",
	// 	Category: "学习",
	// })
	// _, err = client.CreateComment(context.TODO(), &pb.CreateCommentRequest{
	// 	PostId:    1,
	// 	TypeId:    2,
	// 	FatherId:  1,
	// 	Content:   "first comment to comment",
	// 	CreatorId: 2,
	// })
	// _, err = client.UpdatePostInfo(context.TODO(), &pb.UpdatePostInfoRequest{
	// 	Id:       1,
	// 	Content:  "",
	// 	Title:    "",
	// 	Category: "娱乐",
	// })

	post, err := client.ListPost(context.TODO(), &pb.ListPostRequest{
		TypeId:   "1",
		Category: "娱乐",
	})

	fmt.Println("post:", post.List[0].Category)

	if err != nil {
		fmt.Println("err: ", err)
	}
}
