package logrest

import (
	"encoding/json"

	"github.com/Sirupsen/logrus"
)

//const defaultTimeFormat string = "2006-01-02T15:04:05.999Z07:00"
const defaultTimeFormat string = "2006-01-02 15:04:05"


type formatter struct {
	additionalFields map[string]string
	fieldRenameMap   map[string]string
}

func newFormatter(additionalFields, fieldRenameMap map[string]string) *formatter {
	return &formatter{
		additionalFields: additionalFields,
		fieldRenameMap:   fieldRenameMap,
	}
}

func (f *formatter) Format(entry *logrus.Entry) ([]byte, error) {
	record := map[string]interface{}{}
	for key, value := range f.additionalFields {
		name := f.fieldName(key)
		record[name] = value
	}

	for key, value := range entry.Data {
		name := f.fieldName(key)

		switch value := value.(type) {
		case error:
			record[name] = value.Error()
		default:
			record[name] = value
		}
	}

	record[f.fieldName("message")] = entry.Message
	record[f.fieldName("logTime")] = entry.Time.Unix()
	record[f.fieldName("level")] = entry.Level.String()

	json, err := json.Marshal(record)
	return json, err
}

func (f *formatter) fieldName(name string) string {
	if desiredName, ok := f.fieldRenameMap[name]; ok {
		return desiredName
	}

	return name
}
