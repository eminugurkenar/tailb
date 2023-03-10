package read

import (
	"compress/gzip"
	"context"

	"io"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type ObjectStorageReader struct {
	client        *s3.S3
	bucket        string
	prefixes      []string
	objectInfos   []*s3.Object
	objectIndex   int64
	currentObject []byte
	readIndex     int64
	ctx           context.Context
}

func NewObjectStorageReader(client *s3.S3, bucket string, prefixes []string, ctx context.Context) *ObjectStorageReader {
	return &ObjectStorageReader{
		client:        client,
		bucket:        bucket,
		prefixes:      prefixes,
		readIndex:     0,
		objectIndex:   0,
		currentObject: nil,
		ctx:           ctx,
	}
}

func (r *ObjectStorageReader) Read(b []byte) (n int, err error) {
	//fill current object list
	if len(r.objectInfos) == 0 {
		r.objectInfos, err = r.listBucket(r.ctx)
		if err != nil {
			return 0, err
		}
	}

	//current object is empty or reading current object is done, go next object
	if r.readIndex >= int64(len(r.currentObject)) {
		r.objectIndex++
		r.readIndex = 0

		object, err := r.client.GetObjectWithContext(r.ctx, &s3.GetObjectInput{
			Bucket: aws.String(r.bucket),
			Key:    r.objectInfos[r.objectIndex].Key,
		})

		defer object.Body.Close()

		zw, err := gzip.NewReader(object.Body)
		if err != nil {
			return 0, err
		}

		r.currentObject, err = ioutil.ReadAll(zw)
		if err != nil {
			return 0, err
		}
	}

	n = copy(b, r.currentObject[r.readIndex:])
	r.readIndex += int64(n)

	//io.EOF check
	if int64(len(r.objectInfos)) == r.objectIndex && r.readIndex == 0 {
		return n, io.EOF
	}

	return n, nil
}

func (r *ObjectStorageReader) listBucket(ctx context.Context) ([]*s3.Object, error) {
	var objects []*s3.Object

	for _, prefix := range r.prefixes {
		err := r.client.ListObjectsPagesWithContext(ctx, &s3.ListObjectsInput{
			Bucket: aws.String(r.bucket),
			Prefix: aws.String(prefix),
		}, func(page *s3.ListObjectsOutput, lastPage bool) bool {
			for _, object := range page.Contents {
				objects = append(objects, object)
			}
			return lastPage
		})

		if err != nil {
			return nil, err
		}
	}

	return objects, nil
}
