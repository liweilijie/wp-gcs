package gcs

import (
	"cloud.google.com/go/iam"
	"cloud.google.com/go/iam/apiv1/iampb"
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	logging "github.com/ipfs/go-log/v2"
	"google.golang.org/api/option"
	"io"
	"os"
	"time"
)

var log = logging.Logger("gcs")

type Gcs interface {
	Upload(ctx context.Context, params *UploadParams) error
	UploadFile(ctx context.Context, localPath string, bucketPath string) error
	GenerateSignedURL(params *GenerateSignedURLParams) (string, error)
	SetBucketPublicIAM(params *UploadParams) error
}

type gcs struct {
	StorageClient *storage.Client
	ProjectId     string
	Bucket        string
}

type BucketAndObject struct {
	Bucket string
	Object string
}

type UploadParams struct {
	LocalObjPath  string
	UploadObjPath string
	*BucketAndObject
}

type GenerateSignedURLParams struct {
	*BucketAndObject
	ExpirationTime time.Time
	UploadObjPath  string
}

func NewGCS(ctx context.Context, projectId string, bucket string) Gcs {
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(projectId))
	if err != nil {
		log.Fatal("failed to create GCS client", err)
	}
	return &gcs{
		StorageClient: client,
		ProjectId:     projectId,
		Bucket:        bucket,
	}
}

func (g gcs) UploadFile(ctx context.Context, localPath string, bucketPath string) error {
	// open the local file that is intended to be uploaded to GCS.
	// ensure the file is closed at the end.
	file, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("os.Open: %w", err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Errorf("failed to close file: %v", err)
		}
	}(file)

	// set the timeout to 1 minute
	ctx, cancel := context.WithTimeout(ctx, time.Minute*30)
	defer cancel()

	object := g.StorageClient.
		Bucket(g.Bucket).
		Object(bucketPath)

	// set a generation-match precondition to avoid potential race
	// conditions and data corruptions. The request to upload is aborted if the
	// object's generation number does not match your precondition.
	object = object.If(storage.Conditions{DoesNotExist: true})

	// Upload an object with storage.Writer, and close it at the end.
	writer := object.NewWriter(ctx)
	defer func(w *storage.Writer) {
		err := w.Close()
		if err != nil {
			log.Errorf("failed to close writer: %v", err)
		}
	}(writer)

	_, err = io.Copy(writer, file)
	if err != nil {
		return err
	}

	//log.Infof("File %s successfully uploaded to %s.\n", localPath, bucketPath)

	return nil
}

func (g gcs) Upload(ctx context.Context, params *UploadParams) error {
	// open the local file that is intended to be uploaded to GCS.
	// ensure the file is closed at the end.
	file, err := os.Open(params.LocalObjPath + "/" + params.Object)
	if err != nil {
		return fmt.Errorf("os.Open: %w", err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Errorf("failed to close file: %v", err)
		}
	}(file)

	// set the timeout to 1 minute
	ctx, cancel := context.WithTimeout(ctx, time.Minute*1)
	defer cancel()

	object := g.StorageClient.
		Bucket(params.Bucket).
		Object(params.UploadObjPath + "/" + params.Object)

	// set a generation-match precondition to avoid potential race
	// conditions and data corruptions. The request to upload is aborted if the
	// object's generation number does not match your precondition.
	object = object.If(storage.Conditions{DoesNotExist: true})

	// Upload an object with storage.Writer, and close it at the end.
	writer := object.NewWriter(ctx)
	defer func(w *storage.Writer) {
		err := w.Close()
		if err != nil {
			log.Errorf("failed to close writer: %v", err)
		}
	}(writer)

	_, err = io.Copy(writer, file)
	if err != nil {
		return err
	}

	log.Infof("File %s successfully uploaded to %s/%s .\n",
		params.Object, params.Bucket, params.UploadObjPath)

	return nil
}

func (g gcs) GenerateSignedURL(params *GenerateSignedURLParams) (string, error) {
	// Set up the signed URL options.
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: params.ExpirationTime,
	}

	// Signing a URL requires credentials authorized to sign a URL. You can pass
	// these in through SignedURLOptions with one of the following options:
	//    a. a Google service account private key, obtainable from the Google Developers Console
	//    b. a Google Access ID with iam.serviceAccounts.signBlob permissions
	//    c. a SignBytes function implementing custom signing.
	// In this example, none of these options are used, which means the SignedURL
	// function attempts to use the same authentication that was used to instantiate
	// the Storage client. This authentication must include a private key or have
	// iam.serviceAccounts.signBlob permissions.
	signedUrl, err := g.StorageClient.
		Bucket(params.Bucket).
		SignedURL(params.UploadObjPath+"/"+params.Object, opts)
	if err != nil {
		log.Fatalf("Failed to generate signed URL: %v", err)
	}

	log.Info("SignedURL generated successfully")
	return signedUrl, nil
}

func (g gcs) SetBucketPublicIAM(params *UploadParams) error {
	ctx := context.Background()

	policy, err := g.StorageClient.Bucket(params.Bucket).IAM().V3().Policy(ctx)
	if err != nil {
		return fmt.Errorf("Bucket(%q).IAM().V3().Policy: %w", params.Bucket, err)
	}
	role := "roles/storage.objectViewer"
	policy.Bindings = append(policy.Bindings, &iampb.Binding{
		Role:    role,
		Members: []string{iam.AllUsers},
	})
	if err := g.StorageClient.Bucket(params.Bucket).IAM().V3().SetPolicy(ctx, policy); err != nil {
		return fmt.Errorf("Bucket(%q).IAM().SetPolicy: %w", params.Bucket, err)
	}
	log.Infof("Bucket %v is now publicly readable\n", params.Bucket)
	return nil
}
