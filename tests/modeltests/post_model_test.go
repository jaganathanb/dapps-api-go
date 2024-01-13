package modeltests

import (
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/jaganathanb/dapps-api-go/api/models"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gopkg.in/go-playground/assert.v1"
)

func TestFindAllPosts(t *testing.T) {

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table %v\n", err)
	}
	_, _, err = seedUsersAndPosts()
	if err != nil {
		log.Fatalf("Error seeding user and post table %v\n", err)
	}
	posts, err := postInstance.FindAllGSTs(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the posts: %v\n", err)
		return
	}
	assert.Equal(t, len(*posts), 2)
}

func TestSavePost(t *testing.T) {

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatalf("Error user and post refreshing table %v\n", err)
	}

	_, err = seedOneUser()
	if err != nil {
		log.Fatalf("Cannot seed user %v\n", err)
	}

	newPost := models.GST{
		ID:        uuid.Nil,
		GSTIN:     "This is the title",
		TradeName: "This is the content",
		OwnerName: "Test",
	}
	savedPost, err := newPost.SaveGST(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the post: %v\n", err)
		return
	}
	assert.Equal(t, newPost.ID, savedPost.ID)
	assert.Equal(t, newPost.GSTIN, savedPost.GSTIN)
	assert.Equal(t, newPost.TradeName, savedPost.TradeName)
	assert.Equal(t, newPost.OwnerName, savedPost.OwnerName)

}

func TestGetPostByID(t *testing.T) {

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table: %v\n", err)
	}
	post, err := seedOneUserAndOnePost()
	if err != nil {
		log.Fatalf("Error Seeding User and post table")
	}
	foundPost, err := postInstance.FindGSTByID(server.DB, post.GSTIN)
	if err != nil {
		t.Errorf("this is the error getting one user: %v\n", err)
		return
	}
	assert.Equal(t, foundPost.ID, post.ID)
	assert.Equal(t, foundPost.GSTIN, post.GSTIN)
	assert.Equal(t, foundPost.TradeName, post.TradeName)
}

func TestUpdateAPost(t *testing.T) {

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table: %v\n", err)
	}
	post, err := seedOneUserAndOnePost()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}
	postUpdate := models.GST{
		ID:        uuid.Nil,
		GSTIN:     "modiUpdate",
		TradeName: "modiupdate@gmail.com",
		OwnerName: post.OwnerName,
	}
	updatedPost, err := postUpdate.UpdateAGST(server.DB)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}
	assert.Equal(t, updatedPost.ID, postUpdate.ID)
	assert.Equal(t, updatedPost.GSTIN, postUpdate.GSTIN)
	assert.Equal(t, updatedPost.TradeName, postUpdate.TradeName)
	assert.Equal(t, updatedPost.OwnerName, postUpdate.OwnerName)
}

func TestDeleteAPost(t *testing.T) {

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table: %v\n", err)
	}
	post, err := seedOneUserAndOnePost()
	if err != nil {
		log.Fatalf("Error Seeding tables")
	}
	isDeleted, err := postInstance.DeleteAGST(server.DB, post.GSTIN)
	if err != nil {
		t.Errorf("this is the error deleting the user: %v\n", err)
		return
	}
	//one shows that the record has been deleted or:
	// assert.Equal(t, int(isDeleted), 1)

	//Can be done this way too
	assert.Equal(t, isDeleted, int64(1))
}
