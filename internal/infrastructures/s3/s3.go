package s3

import (
	"fmt"
	"refit_backend/internal/logger"

	"github.com/minio/minio-go"
	"github.com/spf13/viper"
)

var s3Client *minio.Client

// GetS3Client ...
func GetS3Client() *minio.Client {
	if s3Client == nil {
		Init()
	}
	return s3Client
}

// Init S3 Client
func Init() {
	accessKey := viper.GetString("digital_ocean.spaces.key")
	secKey := viper.GetString("digital_ocean.spaces.secret")
	endpoint := viper.GetString("digital_ocean.spaces.endpoint")
	ssl := viper.GetBool("digital_ocean.spaces.ssl")

	// Initiate a client using DigitalOcean Spaces.
	client, err := minio.New(endpoint, accessKey, secKey, ssl)
	if err != nil {
		logger.Errorf("could not initiate connection to s3: %s", err.Error())
	}

	s3Client = client
	// test()
}

func test() {

	println("start")
	// Create a done channel to control 'ListObjects' go routine.
	doneCh := make(chan struct{})
	defer close(doneCh)

	// List all objects from a bucket-name with a matching prefix.

	for object := range s3Client.ListObjectsV2("static-luqmanul", "refit/users/1/profile-image", true, doneCh) {
		if object.Err != nil {
			fmt.Println(object.Err)
			return
		}
		fmt.Println("key", object.Key)
		fmt.Println("content-type", object.ContentType)
		fmt.Println("owner", object.Owner.DisplayName, object.Owner.ID)

		objSpec, err := s3Client.GetObject("static-luqmanul", object.Key, minio.GetObjectOptions{})
		if err != nil {
			fmt.Println(err.Error())
		}
		objInfo, err := objSpec.Stat()
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(objInfo.ContentType)
	}

	println("done")

}
