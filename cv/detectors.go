package cv

import (
	"github.com/orsinium-labs/wypes"
	"gocv.io/x/gocv"
)

const maxIndex = 1048560

// CascadeClassifier is a wrapper around gocv.CascadeClassifier for detection.
type CascadeClassifier struct {
	ID       wypes.UInt32
	Name     string
	Filename string

	Classifier gocv.CascadeClassifier
}

// NewCascadeClassifier creates a new CascadeClassifier.
func NewCascadeClassifier(name string) *CascadeClassifier {
	return &CascadeClassifier{
		Name:       name,
		Classifier: gocv.NewCascadeClassifier()}
}

// SetClassifier sets the gocv.CascadeClassifier for the CascadeClassifier.
func (cc *CascadeClassifier) SetClassifier(c gocv.CascadeClassifier) {
	cc.Classifier = c
}

// Close closes the CascadeClassifier.
func (c *CascadeClassifier) Close() {
	c.Classifier.Close()
}

// FaceDetectorYN is a wrapper around gocv.FaceDectectorYN for detection.
type FaceDetectorYN struct {
	ID       wypes.UInt32
	Name     string
	Filename string

	Detector gocv.FaceDetectorYN
}

// NewFaceDetectorYN creates a new FaceDetectorYN.
func NewFaceDetectorYN(model string) *FaceDetectorYN {
	return &FaceDetectorYN{
		Filename: model,
	}
}

// Close closes the FaceDectectorYN.
func (d *FaceDetectorYN) Close() {
	d.Detector.Close()
}

// SetDetector sets the gocv.FaceDetectorYN for the FaceDectectorYN.
func (d *FaceDetectorYN) SetDetector(dd gocv.FaceDetectorYN) {
	d.Detector = dd
}
