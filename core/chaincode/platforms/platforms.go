package platforms

import (
	"archive/tar"
	"fmt"
	"github.com/idci/core/chaincode/platforms/golang"
	pb "github.com/idci/protos"
)

type Platform interface {
	ValidateSpec(spec *pb.ChaincodeSpec) error
	WritePackage(spec *pb.ChaincodeSpec, tw *tar.Writer) error
}

func Find(chaincodeType pb.ChaincodeSpec_Type) (Platform, error) {

	switch chaincodeType {
	case pb.ChaincodeSpec_GOLANG:
		return &golang.Platform{}, nil
	default:
		return nil, fmt.Errorf("Unknown chaincodeType: %s", chaincodeType)
	}

}
