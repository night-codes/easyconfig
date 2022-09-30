package easyconfig

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/fatih/color"
	"github.com/iancoleman/strcase"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
)

type (
	Loader struct {
		Sources        []Source
		HelpMSG        string
		DisableHelpMsg bool
	}

	errCollector []error

	Source interface {
		Load(structPtr interface{}) error
	}

	// FileSource satisifies the loader interface. It loads the
	// configuration from the given file.
	FileSource struct {
		Path string
	}

	// FileSource loads configuration from the given .json file.
	JSONSource struct {
		Path string
	}

	// YAMLSource loads configuration from the given .yaml or .yml file.
	YAMLSource struct {
		Path string
	}

	// TOMLSource loads configuration from the given .toml file.
	TOMLSource struct {
		Path string
	}

	// DirSource loads configuration from kubernetes ConfigMap and Secret mountPath directory
	DirSource struct {
		Path string
	}

	// EnvFileSource loads configuration from the given .env file.
	EnvFileSource struct {
		Prefix string
		Path   string
	}

	// EnvSource loads configuration from environment variables
	EnvSource struct {
		Prefix string
	}

	// FlagsSource loads configuration from flags
	FlagsSource struct {
	}

	strErr string
)

var (
	Acronims = []string{"API", "SMTP", "PostgreSQL", "SQL", "JSON", "YAML", "DB", "AI", "CRM", "HTTPS", "HTTP", "FTP", "SSH"}
	bold     = color.New(color.Bold)
)

func (e strErr) Error() string {
	return string(e)
}

const (
	ErrIsDirectory     strErr = "file is a directory"
	ErrUnknownFileType strErr = "unknown file type"
)

func (c *errCollector) Collect(e error) {
	if e != nil {
		*c = append(*c, e)
	}
}

func (c *errCollector) Error() (err error) {
	if len(*c) == 0 {
		return nil
	}
	ret := ""
	for _, e := range *c {
		ret += fmt.Sprintf("%s\n", e.Error())
	}

	return errors.New(ret)
}

func NewLoader(sources []Source, helpMessage ...string) *Loader {
	for _, acron := range Acronims {
		strcase.ConfigureAcronym(acron, strings.ToLower(acron))
	}
	msg := fmt.Sprintf("\nUsage:\n    %s [arguments]", filepath.Base(os.Args[0]))

	if len(helpMessage) > 0 {
		msg = strings.Join(helpMessage, "\n")
	}
	return &Loader{
		Sources:        sources,
		HelpMSG:        msg,
		DisableHelpMsg: len(msg) == 0,
	}
}

// Load configuration
func (l Loader) Load(structPtr interface{}) error {
	if !l.DisableHelpMsg && len(os.Args) == 2 && strings.HasSuffix(strings.ToLower(strings.Split(os.Args[1], "=")[0]), "help") {
		l.Help(structPtr)
		os.Exit(0)
	}
	errs := new(errCollector)
	for _, src := range l.Sources {
		errs.Collect(src.Load(structPtr))
	}
	return errs.Error()
}

// Load configuration
func (l Loader) Help(structPtr interface{}) {
	structElem := reflect.ValueOf(structPtr).Elem()
	tags := []string{}
	fromEnv := false
	fromFlag := false
	fromDir := false
	prefix := ""
	for _, source := range l.Sources {
		if t, ok := source.(EnvFileSource); ok {
			fromEnv = true
			prefix = t.Prefix
		}
		if t, ok := source.(EnvSource); ok {
			fromEnv = true
			prefix = t.Prefix
		}
		if t, ok := source.(*EnvFileSource); ok {
			fromEnv = true
			prefix = t.Prefix
		}
		if t, ok := source.(*EnvSource); ok {
			prefix = t.Prefix
			fromEnv = true
		}
		if _, ok := source.(FlagsSource); ok {
			fromFlag = true
		}
		if _, ok := source.(*FlagsSource); ok {
			fromFlag = true
		}
		if _, ok := source.(DirSource); ok {
			fromDir = true
		}
		if _, ok := source.(*DirSource); ok {
			fromDir = true
		}
		if _, ok := source.(FileSource); ok {
			fromDir = true
		}
		if _, ok := source.(*FileSource); ok {
			fromDir = true
		}
	}
	if fromFlag {
		tags = append(tags, "flag")
	}
	if fromEnv {
		tags = append(tags, "env")
	}
	if fromDir {
		tags = append(tags, "dir")
	}

	fmt.Println(l.HelpMSG)

	for _, tag := range tags {
		switch tag {
		case "flag":
			fmt.Printf("\nThe commands are:\n\n")
		case "dir":
			fmt.Printf("\nConfiguration directory files to use:\n\n")
		default:
			fmt.Printf("\nEnvironment variables to use:\n\n")
		}
		for i := 0; i < structElem.NumField(); i++ {
			field := structElem.Type().Field(i)
			elem := structElem.Field(i).Interface()
			tagVal := field.Tag.Get(tag)
			if tagVal == "-" {
				continue
			}
			fieldName, _ := convertName(field.Name, tag, tagVal, prefix)

			bold.Printf("    %s\n", fieldName)
			if tag == "flag" {
				printed := false
				switch v := elem.(type) {
				case string:
					if v != "" {
						printed = true
						fmt.Printf("        Set value of %s. Default: %q\n", field.Name, v)
					}
				case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
					if v != 0 {
						printed = true
						fmt.Printf("        Set number value of %s. Default: %v\n", field.Name, elem)
					}
				case bool:
					printed = true
					fmt.Printf("        Set value of %s. Default: %v\n", field.Name, elem)
				}

				if !printed {
					fmt.Printf("        Set value of %s\n", field.Name)
				}
			}
		}
	}
}

// Load config from file
func (s FileSource) Load(structPtr interface{}) error {
	if info, err := os.Stat(s.Path); err != nil {
		return err
	} else {
		if info.IsDir() {
			return (&DirSource{Path: s.Path}).Load(structPtr)
		}
		ext := filepath.Ext(s.Path)
		switch ext {
		case ".json":
			return (&JSONSource{Path: s.Path}).Load(structPtr)
		case ".yaml", ".yml":
			return (&YAMLSource{Path: s.Path}).Load(structPtr)
		case ".env":
			return (&EnvFileSource{Path: s.Path}).Load(structPtr)
		case ".toml":
			return (&TOMLSource{Path: s.Path}).Load(structPtr)
		default:
			if err := (&JSONSource{Path: s.Path}).Load(structPtr); err == nil {
				return nil
			}
			if err := (&YAMLSource{Path: s.Path}).Load(structPtr); err == nil {
				return nil
			}
			if err := (&EnvFileSource{Path: s.Path}).Load(structPtr); err == nil {
				return nil
			}
			if err := (&TOMLSource{Path: s.Path}).Load(structPtr); err == nil {
				return nil
			}
			return ErrUnknownFileType
		}
	}
}

// Load JSON configuration file
func (s JSONSource) Load(structPtr interface{}) error {
	file, err := getFile(s.Path)
	if err != nil {
		return err
	}
	return json.NewDecoder(file).Decode(structPtr)
}

// Load YAML configuration file
func (s YAMLSource) Load(structPtr interface{}) error {
	file, err := getFile(s.Path)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, structPtr)
}

// Load TOML configuration file
func (s TOMLSource) Load(structPtr interface{}) error {
	file, err := getFile(s.Path)
	if err != nil {
		return err
	}
	_, err = toml.DecodeReader(file, structPtr)
	return err
}

// Load ENV configuration file
func (s EnvFileSource) Load(structPtr interface{}) error {
	file, err := getFile(s.Path)
	if err != nil {
		return err
	}

	envMap, err := godotenv.Parse(file)
	if err != nil {
		return err
	}

	return map2struct("env", s.Prefix, envMap, structPtr)
}

// Load configuration from environment variables
func (s EnvSource) Load(structPtr interface{}) error {
	envMap := map[string]string{}
	for _, s := range os.Environ() {
		if strings.Contains(s, "=") {
			p := strings.SplitN(s, "=", 2)
			envMap[p[0]] = p[1]
		}
	}
	return map2struct("env", s.Prefix, envMap, structPtr)
}

// Load configuration from environment variables
func (s DirSource) Load(structPtr interface{}) error {
	dirMap := map[string]string{}
	files, err := ioutil.ReadDir(s.Path)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() && file.Size() < 10485760 { // 10 Mb
			data, err := ioutil.ReadFile(filepath.Join(s.Path, file.Name()))
			if err != nil {
				continue
			}
			dirMap[file.Name()] = string(data)
		}
	}
	return map2struct("dir", "", dirMap, structPtr)
}

//  Load configuration from the command-line.
func (s FlagsSource) Load(structPtr interface{}) error {
	argsMap := map[string]string{}
	prev := ""
	for _, s := range os.Args[1:] {
		if strings.HasPrefix(s, "-") {
			prev = ""
			if strings.Contains(s, "=") {
				p := strings.SplitN(s, "=", 2)
				argsMap[p[0]] = p[1]
			} else {
				prev = s
			}
		} else if prev != "" {
			argsMap[prev] = s
		}
	}
	return map2struct("flag", "", argsMap, structPtr)
}

func getFile(path string) (*os.File, error) {
	pwd, err := os.Getwd()
	if !filepath.IsAbs(path) && err == nil {
		path = filepath.Join(pwd, path)
	}
	if _, err := os.Stat(path); err != nil {
		return nil, err
	}
	return os.Open(path)
}

func map2struct(tag, prefix string, mp map[string]string, structPtr interface{}) error {
	for key, value := range mp {
		key = strings.TrimSpace(key)
		if key != "" && value != "" {
			if structPtr != nil {
				setValue(tag, prefix, structPtr, key, value)
			}
		}
	}

	return nil
}

func setValue(tag, prefix string, structPtr interface{}, key string, value string) {

	structElem := reflect.ValueOf(structPtr).Elem()

	for i := 0; i < structElem.NumField(); i++ {
		field := structElem.Type().Field(i)
		tagVal := strings.TrimSpace(field.Tag.Get(tag))
		if tagVal == "-" {
			continue
		}
		fieldName, separator := convertName(field.Name, tag, tagVal, prefix)
		if tag == "env" {
			key = strings.ToUpper(key)
		}
		elemPtr := structElem.Field(i).Addr().Interface()

		if fieldName == key {
			switch field.Type.Kind() {
			case reflect.Slice:
				setSlice(elemPtr, value, separator)
			default:
				setField(elemPtr, value)
			}
		}
	}
}

func convertName(name, tag, tagVal, prefix string) (fieldName, separator string) {
	separator = ","
	if tag == "env" {
		separator = ":"
		fieldName = strcase.ToScreamingSnake(name)
	} else if tag == "dir" {
		fieldName = strcase.ToKebab(name)
	} else {
		fieldName = name
		for _, acron := range Acronims {
			if strings.HasPrefix(fieldName, acron) {
				fieldName = strings.Replace(fieldName, acron, strings.ToLower(acron), 1)
			}
		}
		fieldName = "-" + strcase.ToLowerCamel(fieldName)
	}

	if strings.Contains(tagVal, ",") {
		p := strings.SplitN(tagVal, ",", 2)
		tagVal = p[0]
		separator = p[1]
	}

	if tagVal != "" {
		fieldName = tagVal
		if tag == "flag" && strings.HasPrefix(fieldName, "-") {
			fieldName = "-" + fieldName
		}
	}

	if tag == "env" {
		prefix = strings.TrimSpace(prefix)
		if prefix != "" && !strings.HasSuffix(prefix, "_") && !strings.HasSuffix(prefix, "-") {
			prefix += "_"
		}
		if prefix != "" && !strings.HasPrefix(fieldName, prefix) {
			fieldName = prefix + fieldName
		}
		fieldName = strings.ToUpper(fieldName)
	}

	return fieldName, separator
}

func setSlice(slicePtr interface{}, value string, separator string) {
	if separator == "" {
		separator = ":"
	}
	slice := reflect.ValueOf(slicePtr).Elem()
	for _, part := range strings.Split(value, separator) {
		fieldPtr := reflect.New(slice.Type().Elem())
		setField(fieldPtr.Interface(), part)
		slice.Set(reflect.Append(slice, fieldPtr.Elem()))
	}
}

func setField(fieldPtr interface{}, fieldValue string) {
	field := reflect.ValueOf(fieldPtr).Elem()
	if field.IsValid() && field.CanSet() {
		switch field.Kind() {
		case reflect.String:
			field.SetString(fieldValue)
		case reflect.Int:
			r, _ := strconv.ParseInt(fieldValue, 10, 0)
			field.SetInt(r)
		case reflect.Int8:
			r, _ := strconv.ParseInt(fieldValue, 10, 8)
			field.SetInt(r)
		case reflect.Int16:
			r, _ := strconv.ParseInt(fieldValue, 10, 16)
			field.SetInt(r)
		case reflect.Int32:
			r, _ := strconv.ParseInt(fieldValue, 10, 32)
			field.SetInt(r)
		case reflect.Int64:
			r, _ := strconv.ParseInt(fieldValue, 10, 64)
			field.SetInt(r)
		case reflect.Uint:
			r, _ := strconv.ParseUint(fieldValue, 10, 0)
			field.SetUint(r)
		case reflect.Uint8:
			r, _ := strconv.ParseUint(fieldValue, 10, 8)
			field.SetUint(r)
		case reflect.Uint16:
			r, _ := strconv.ParseUint(fieldValue, 10, 16)
			field.SetUint(r)
		case reflect.Uint32:
			r, _ := strconv.ParseUint(fieldValue, 10, 32)
			field.SetUint(r)
		case reflect.Uint64:
			r, _ := strconv.ParseUint(fieldValue, 10, 64)
			field.SetUint(r)
		case reflect.Bool:
			r, _ := strconv.ParseBool(fieldValue)
			field.SetBool(r)
		case reflect.Float32:
			r, _ := strconv.ParseFloat(fieldValue, 32)
			field.SetFloat(r)
		case reflect.Float64:
			r, _ := strconv.ParseFloat(fieldValue, 64)
			field.SetFloat(r)
		}
	}
}
