package logging

import "github.com/sirupsen/logrus"

//New : instantiate new logging
func New(packageName string, functionName string) *logrus.Entry {
	log := logrus.WithField("package", packageName).WithField("function", functionName)
	return log
}
