package baiducloud

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/baidubce/bce-sdk-go/bce"
)

// A default message of ComplexError's Err. It is format to Resource <resource-id> <operation> Failed!!! <error source>
const DefaultErrorMsg = "Resource %s %s Failed!!! %s"

// ComplexError is a format error which inclouding origin error, extra error message, error occurred file and line
// Cause: a error is a origin error that comes from SDK, some expections and so on
// Err: a new error is built from extra message
// Path: the file path of error occurred
// Line: the file line of error occurred
type ComplexError struct {
	Cause error
	Err   error
	Path  string
	Line  int
}

type ErrorSource string

const (
	// common error
	ResourceNotFound          = "ResourceNotfound"
	ResourceNotFound2         = "ResourceNotFound"
	ResourceNotFoundException = "ResourceNotFoundException"

	// bcc error
	OperationDenied      = "OperationDenied"
	ReleaseWhileCreating = "Instance.ReleaseWhileCreating"

	// not found error
	InstanceNotFound = "InstanceNotFound"
	EipNotFound      = "EipNotFound"
	CceNotFound      = "Cce.warning.ClusterNotExist"
	InstanceNotExist = "InstanceNotExist"

	// scs error
	InvalidInstanceStatus = "InvalidInstanceStatus"
	OperationException    = "OperationException"
)

const (
	// sdk error
	BCESDKGoERROR = ErrorSource("[SDK bce-sdk-go ERROR]")
)

const GetFailTargetStatus = "Get Fail target status: %s."

var (
	// bcc error
	BccNotFound = []string{"InvalidInstanceId.NotFound", "Forbidden.InstanceNotFound"}

	// nat gateway error
	NatGatewayNotFound = []string{"NoSuchNat"}

	// peer conn error
	PeerConnNotFound = []string{"EOF"}

	// common error
	ObjectNotFound = []string{"NoSuchObject"}

	// replication configuration error
	ReplicationConfigurationNotFound = []string{"NoReplicationConfiguration"}

	// cce error
	CceClusterNotFound = []string{CceNotFound}
)

var NotFoundErrorList = []string{
	ResourceNotFoundException, ResourceNotFound, ResourceNotFound2, InstanceNotFound,
	"NoSuchObject", "NoSuchNat", EipNotFound, InstanceNotExist, CceNotFound,
}

// An Error to wrap the different erros
type WrapErrorOld struct {
	originError error
	errorSource ErrorSource
	errorPath   string
	message     string
	suggestion  string
}

func (e *WrapErrorOld) Error() string {
	return fmt.Sprintf("[ERROR] %s: %s %s:\n%s\n%s", e.errorPath, e.message, e.errorSource, e.originError.Error(), e.suggestion)
}

func NotFoundError(err error) bool {
	if e, ok := err.(*WrapErrorOld); ok {
		err = e.originError
	}
	if err == nil {
		return false
	}
	if e, ok := err.(*ComplexError); ok {
		if e.Err != nil {
			for _, notFoundErr := range NotFoundErrorList {
				if strings.HasPrefix(e.Err.Error(), notFoundErr) {
					return true
				}
			}
		}
		return NotFoundError(e.Cause)
	}

	if e, ok := err.(*bce.BceServiceError); ok {
		if stringInSlice(NotFoundErrorList, e.Code) {
			return true
		}
	}

	for _, notFoundErr := range NotFoundErrorList {
		if strings.HasPrefix(err.Error(), notFoundErr) {
			return true
		}
	}

	return false
}

func IsExceptedErrors(err error, expectCodes []string) bool {
	if e, ok := err.(*WrapErrorOld); ok {
		err = e.originError
	}
	if err == nil {
		return false
	}

	if e, ok := err.(*ComplexError); ok {
		return IsExceptedErrors(e.Cause, expectCodes)
	}

	for _, code := range expectCodes {
		if e, ok := err.(bce.BceError); ok && (strings.Contains(e.Error(), "Code: "+code)) {
			return true
		}
		if strings.Contains(err.Error(), code) {
			return true
		}
	}
	return false
}

func (e ComplexError) Error() string {
	if e.Cause == nil {
		e.Cause = Error("<nil cause>")
	}
	if e.Err == nil {
		return fmt.Sprintf("[ERROR] %s:%d:\n%s", e.Path, e.Line, e.Cause.Error())
	}
	return fmt.Sprintf("[ERROR] %s:%d: %s:\n%s", e.Path, e.Line, e.Err.Error(), e.Cause.Error())
}

func Error(msg string, args ...interface{}) error {
	return fmt.Errorf(msg, args...)
}

func WrapComplexError(cause, err error, filepath string, fileline int) error {
	return &ComplexError{
		Cause: cause,
		Err:   err,
		Path:  filepath,
		Line:  fileline,
	}
}

// Return a ComplexError which including error occurred file and path
func WrapError(cause error) error {
	if cause == nil {
		return nil
	}
	_, filepath, line, ok := runtime.Caller(1)
	if !ok {
		log.Printf("[ERROR] runtime.Caller error in WrapError.")
		return WrapComplexError(cause, nil, "", -1)
	}
	parts := strings.Split(filepath, "/")
	if len(parts) > 3 {
		filepath = strings.Join(parts[len(parts)-3:], "/")
	}
	return WrapComplexError(cause, nil, filepath, line)
}

// Return a ComplexError which including extra error message, error occurred file and path
func WrapErrorf(cause error, msg string, args ...interface{}) error {
	if cause == nil && strings.TrimSpace(msg) == "" {
		return nil
	}
	_, filepath, line, ok := runtime.Caller(1)
	if !ok {
		log.Printf("[ERROR] runtime.Caller error in WrapErrorf.")
		return WrapComplexError(cause, Error(msg), "", -1)
	}
	parts := strings.Split(filepath, "/")
	if len(parts) > 3 {
		filepath = strings.Join(parts[len(parts)-3:], "/")
	}
	return WrapComplexError(cause, fmt.Errorf(msg, args...), filepath, line)
}
