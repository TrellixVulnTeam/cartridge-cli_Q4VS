package common

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	goVersion "github.com/hashicorp/go-version"
	"github.com/robfig/config"
	"github.com/tarantool/cartridge-cli/cli/connector"
)

var (
	// Used for shallow validation of a version string.
	tarantoolVersionRegexp *regexp.Regexp
	// Used for thorough validation of a version string and to extract string version to a struct.
	tarantoolVersionFullRegexp *regexp.Regexp
)

const (
	tarantoolExeName         = "tarantool"
	tarantoolVersionFlag     = "--version"
	TarantoolVersionFileName = "tarantool.txt"
	tarantoolVersionOptName  = "TARANTOOL"
)

var tarantoolVersionShortRegexps = []*regexp.Regexp{
	regexp.MustCompile(`^(?P<Major>\d+)$`),
	regexp.MustCompile(`^(?P<Major>\d+)\.(?P<Minor>\d+)$`),
	regexp.MustCompile(`^(?P<Major>\d+)\.(?P<Minor>\d+)\.(?P<Patch>\d+)$`),
	regexp.MustCompile(`^(?P<Major>\d+)\.(?P<Minor>\d+)\.(?P<Patch>\d+)` +
		`[-~](?P<TagSuffix>alpha\d+|beta\d+|rc\d+|entrypoint)$`),
}

func init() {
	tarantoolVersionRegexp = regexp.MustCompile(`\d+\.\d+\.\d+[-\w]*`)
	// Part 1 is semVer X.Y.Z ,
	// part 2 (optional) is a tag suffix for pre-release,
	// part 3 is number of commits since tag and commit hash,
	// part 4 (optional) is enterprise suffix,
	// part 5 (optional) is development build suffix.
	tarantoolVersionFullRegexp = regexp.MustCompile(
		`^(?P<Major>\d+)\.(?P<Minor>\d+)?\.(?P<Patch>\d+)` +
			`(?:-(?P<TagSuffix>alpha\d+|beta\d+|rc\d+|entrypoint))?` +
			`-(?P<CommitsSinceTag>\d+)-(?P<CommitHashId>g[0-9a-f]+)` +
			`(?:-(?P<EnterpriseSDKRevision>r\d+)(?:-(?P<EnterpriseIsOnMacOS>macos))?)?` +
			`(?:-(?P<IsDevelopmentBuild>dev))?$`)
}

// findNamedMatches processes regexp with named capture groups
// and transforms output to a map. If capture group is optional
// and was not found, map value is empty string.
func findNamedMatches(regex *regexp.Regexp, str string) map[string]string {
	match := regex.FindStringSubmatch(str)

	if match == nil {
		return nil
	}

	results := map[string]string{}
	for i, name := range match {
		// i == 0 is input string.
		if i != 0 {
			results[regex.SubexpNames()[i]] = name
		}
	}
	return results
}

// findNamedMatchesMultipleRegexps processes a slice of regexps with
// named capture groups and transforms output to a map. If capture group
// is optional and was not found, map value is empty string.
func findNamedMatchesMultipleRegexps(regexps []*regexp.Regexp, str string) map[string]string {
	results := map[string]string{}
	for _, regex := range regexps {
		results = findNamedMatches(regex, str)
		if results != nil {
			break
		}
	}

	return results
}

type TarantoolVersion struct {
	Major                 uint64
	Minor                 OptionalUint64
	Patch                 OptionalUint64
	TagSuffix             string
	CommitsSinceTag       uint64
	CommitHashId          string
	EnterpriseSDKRevision string
	EnterpriseIsOnMacOS   bool
	IsDevelopmentBuild    bool
}

func atoiUint64(str string) (uint64, error) {
	res, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("Failed to parse Tarantool version: cast to number error (%s)", err)
	}
	return res, nil
}

// ParseShortTarantoolVersion parse Tarantool version provided by user.
func ParseShortTarantoolVersion(version string) (TarantoolVersion, error) {
	var result TarantoolVersion
	var err error

	matches := findNamedMatchesMultipleRegexps(tarantoolVersionShortRegexps, version)
	if matches == nil || len(matches) == 0 {
		return result, fmt.Errorf("Failed to parse Tarantool version: format is not valid")
	}

	if result.Major, err = atoiUint64(matches["Major"]); err != nil {
		return result, err
	}

	if matches["Minor"] == "" {
		return result, nil
	} else if minor, err := atoiUint64(matches["Minor"]); err == nil {
		result.Minor = OptFrom(minor)
	} else {
		return result, err
	}

	if matches["Patch"] == "" {
		return result, nil
	} else if patch, err := atoiUint64(matches["Patch"]); err == nil {
		result.Patch = OptFrom(patch)
	} else {
		return result, err
	}

	result.TagSuffix = matches["TagSuffix"]

	if result.TagSuffix == "entrypoint" {
		return result, fmt.Errorf("Entrypoint build cannot be used for --tarantool-version")
	}

	return result, nil
}

// ParseTarantoolVersion extracts Tarantool version string to a TarantoolVersion struct.
func ParseTarantoolVersion(version string) (TarantoolVersion, error) {
	var result TarantoolVersion
	var err error

	matches := findNamedMatches(tarantoolVersionFullRegexp, version)
	if matches == nil || len(matches) == 0 {
		return result, fmt.Errorf("Failed to parse Tarantool version: format is not valid")
	}

	if result.Major, err = atoiUint64(matches["Major"]); err != nil {
		return result, err
	}

	if minor, err := atoiUint64(matches["Minor"]); err == nil {
		result.Minor = OptFrom(minor)
	} else {
		return result, err
	}

	if patch, err := atoiUint64(matches["Patch"]); err == nil {
		result.Patch = OptFrom(patch)
	} else {
		return result, err
	}

	result.TagSuffix = matches["TagSuffix"]

	if result.CommitsSinceTag, err = atoiUint64(matches["CommitsSinceTag"]); err != nil {
		return result, err
	}

	result.CommitHashId = matches["CommitHashId"]

	result.EnterpriseSDKRevision = matches["EnterpriseSDKRevision"]

	result.EnterpriseIsOnMacOS = (matches["EnterpriseIsOnMacOS"] != "")

	result.IsDevelopmentBuild = (matches["IsDevelopmentBuild"] != "")

	return result, nil
}

// GetTarantoolVersion gets Tarantool version from the environment.
func GetTarantoolVersion(tarantoolDir string) (string, error) {
	var err error

	tarantool := filepath.Join(tarantoolDir, tarantoolExeName)
	versionCmd := exec.Command(tarantool, tarantoolVersionFlag)

	tarantoolVersion, err := GetOutput(versionCmd, nil)
	if err != nil {
		return "", err
	}

	tarantoolVersion = tarantoolVersionRegexp.FindString(tarantoolVersion)

	if tarantoolVersion == "" {
		return "", fmt.Errorf("Failed to match Tarantool version")
	}

	return tarantoolVersion, nil
}

// GetTarantoolVersionFromFile gets Tarantool version from config file.
func GetTarantoolVersionFromFile(tarantoolVersionFilePath string) (string, error) {
	tarantoolConfig, err := config.ReadDefault(tarantoolVersionFilePath)
	if err != nil {
		return "", fmt.Errorf("%s is specified in bad format: %s", TarantoolVersionFileName, err)
	}

	var tarantoolVersion string

	tarantoolVersion, _ = tarantoolConfig.RawStringDefault(tarantoolVersionOptName)

	if tarantoolVersion == "" {
		return "", fmt.Errorf(
			"You should specify specify %s in %s file",
			tarantoolVersionOptName, TarantoolVersionFileName,
		)
	}

	return tarantoolVersion, nil
}

// TarantoolIsEnterprise checks if Tarantool is Enterprise
func TarantoolIsEnterprise(tarantoolDir string) (bool, error) {
	tarantool := filepath.Join(tarantoolDir, tarantoolExeName)
	versionCmd := exec.Command(tarantool, tarantoolVersionFlag)

	tarantoolVersion, err := GetOutput(versionCmd, nil)
	if err != nil {
		return false, err
	}

	return strings.HasPrefix(tarantoolVersion, "Tarantool Enterprise"), nil
}

// GetTarantoolDir returns Tarantool executable directory.
func GetTarantoolDir() (string, error) {
	var err error

	tarantool, err := exec.LookPath(tarantoolExeName)
	if err != nil {
		return "", fmt.Errorf("tarantool executable not found")
	}

	return filepath.Dir(tarantool), nil
}

// GetMajorMinorVersion computes returns `<major>.<minor>` string
// for a given version
func GetMajorMinorVersion(version string) string {
	parts := strings.SplitN(version, ".", 3)
	major := parts[0]
	minor := parts[1]

	majorMinorVersion := fmt.Sprintf("%s.%s", major, minor)

	return majorMinorVersion
}

// GetMinimalRequiredVersion computes minimal required Tarantool version for a package (rpm, deb).
func GetMinimalRequiredVersion(ver TarantoolVersion) (string, error) {
	// Old-style package version policy allowed X.Y.Z-N versions for N > 0 .
	minor, ok := OptValue(ver.Minor)
	if !ok {
		return "", fmt.Errorf("Version minor part is not specified")
	}
	patch, ok := OptValue(ver.Patch)
	if !ok {
		return "", fmt.Errorf("Version patch part is not specified")
	}

	if (ver.Major == 2 && minor <= 8) || (ver.Major < 2) {
		return fmt.Sprintf("%d.%d.%d.%d", ver.Major, minor, patch, ver.CommitsSinceTag), nil
	}

	if ver.IsDevelopmentBuild {
		return "", fmt.Errorf("Can't compute minimal required version for a development build")
	}

	if ver.TagSuffix == "entrypoint" {
		return "", fmt.Errorf("Can't compute minimal required version for an entrypoint build")
	}

	if ver.TagSuffix != "" {
		return fmt.Sprintf("%d.%d.%d~%s", ver.Major, minor, patch, ver.TagSuffix), nil
	}

	return fmt.Sprintf("%d.%d.%d", ver.Major, minor, patch), nil
}

// GetNextMajorVersion computes next Major version for a given one.
// For example, for 1.10.3 it's 2 .
func GetNextMajorVersion(ver TarantoolVersion) string {
	return strconv.Itoa(int(ver.Major) + 1)
}

func GetCartridgeVersionStr(conn *connector.Conn) (string, error) {
	req := connector.EvalReq(getCartridgeVersionBody).SetReadTimeout(3 * time.Second)

	var versionStrSlice []string
	if err := conn.ExecTyped(req, &versionStrSlice); err != nil {
		return "", fmt.Errorf("Failed to eval get Cartridge version function: %s", err)
	}

	if len(versionStrSlice) != 1 {
		return "", fmt.Errorf("Cartridge version received in a wrong format")
	}

	versionStr := versionStrSlice[0]

	return versionStr, nil
}

func GetMajorCartridgeVersion(conn *connector.Conn) (int, error) {
	cartridgeVersionStr, err := GetCartridgeVersionStr(conn)
	if err != nil {
		return 0, err
	}

	if cartridgeVersionStr == "scm-1" {
		return 2, nil
	}

	cartridgeVersion, err := goVersion.NewSemver(cartridgeVersionStr)
	if err != nil {
		return 0, fmt.Errorf("Failed to parse Tarantool version: %s", err)
	}

	major := cartridgeVersion.Segments()[0]

	return major, nil
}

// FindRockspec finds *.rockspec file in specified path
func FindRockspec(path string) (string, error) {
	rockspecs, err := filepath.Glob(filepath.Join(path, "*.rockspec"))

	if err != nil {
		return "", fmt.Errorf("Failed to find rockspec: %s", err)
	}

	if len(rockspecs) > 0 {
		return rockspecs[0], nil
	}

	return "", nil
}

const (
	getCartridgeVersionBody = `return require('cartridge').VERSION or '1'`
)
