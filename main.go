package main

import (
	"context"
	"fmt"
	bson "github.com/mongodb/mongo-go-driver/bson"
	mongo "github.com/mongodb/mongo-go-driver/mongo"
	"os"
	//	"time"
	"bufio"
	//	"io/ioutil"
)

type Post struct {
	Title   string `bson:"title,omitempty"`
	Date    string `bson:"date,omitempty"`
	Tags    string `bson:"tags,omitempty"`
	Draft   string `bson:"draft,omitempty"`
	Content string `bson:"content,omitempty"`
}

func main() {
	for true {
		content := getContent()
		postDirectory := "/Users/jak/tools/beaglechow/publish/hugo/beagleblog.com/content/post/"
		for _, post := range content {
			fmt.Printf("writing file to %s%s.md\n", string(postDirectory), string(post.Title))
			fmt.Printf("content dump:\n%s\n\n", string(post.Content))
			f, err := os.Create(postDirectory + post.Title + ".md")
			if err != nil {
				panic(err)
			}
			//todo: error test here
			defer f.Close()
			w := bufio.NewWriter(f)
			w.WriteString("---\ntitle: \"" + post.Title + "\"\ndate: " + post.Date + "\ntags: []\ndraft: " + post.Draft + "\n---\n" + post.Content + "\n")
			w.Flush()
		}
	}
}

func getContent() []Post {
	//	mongodbUrl := os.Getenv("MONGODB_URL")
	mongodbUrl := "mongodb://localhost:27017"
	client, err := mongo.Connect(context.Background(), mongodbUrl, nil)
	db := client.Database("beagleblog")
	coll := db.Collection("posts")
	results := []Post{}
	cursor, err := coll.Find(context.Background(), bson.NewDocument(bson.EC.String("draft", "false")))
	for cursor.Next(context.Background()) {
		elem := Post{}
		err := cursor.Decode(&elem)
		if err != nil {
			panic(err)
		}
		results = append(results, elem)
	}
	if err != nil {
		panic(err)
	}
	return results
}
