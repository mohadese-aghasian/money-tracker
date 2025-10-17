package utils

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"

	// "path/filepath"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os/exec"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unicode"
)

func validateNoNumbers(name string) bool {
	// name := fl.Field().String()
	re := regexp.MustCompile(`^[^\d]+$`) // No digits allowed
	return re.MatchString(name)
}

func isValidEmail(email string) bool {
	// A very simple regex for email validation
	re := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	return regexp.MustCompile(re).MatchString(email)
}

func IsEmpty(item interface{}) bool {
	if item == nil {
		return true
	}

	value := reflect.ValueOf(item)
	kind := value.Kind()

	switch kind {
	case reflect.String:
		return value.Len() == 0 || value.String() == "undefined" || value.String() == "NaN" || value.String() == "nan"
	case reflect.Slice, reflect.Array, reflect.Map:
		return value.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return value.IsNil()
	case reflect.Invalid:
		return true
	}

	return false
}

func IsTextAnswerCorrect(answer string, correct_answer string) bool {
	return strings.EqualFold(correct_answer, answer)

}

func GenerateFileName(originalName string) string {
	// ext := filepath.Ext(originalName)  // e.g. ".jpg"
	timestamp := time.Now().UnixNano() // high-res timestamp
	randomPart := rand.Intn(1000000)   // random number up to 999999
	return fmt.Sprintf("%d_%d", timestamp, randomPart)
}

func GenerateSlugUnicode(text string) string {
	// 1. Lowercase (safe for Persian too)
	slug := strings.ToLower(text)

	// 2. Replace spaces and underscores with hyphen
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")

	// 3. Remove disallowed characters (keep Persian, English letters, digits, hyphen)
	// This function keeps Unicode letters and digits
	var b strings.Builder
	for _, r := range slug {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' {
			b.WriteRune(r)
		}
		// else drop punctuation like .,!? etc.
	}
	slug = b.String()

	// 4. Replace multiple hyphens with single hyphen
	reg := regexp.MustCompile(`-+`)
	slug = reg.ReplaceAllString(slug, "-")

	// 5. Trim leading/trailing hyphens
	slug = strings.Trim(slug, "-")

	return slug
}

func GetEnvInt(key string, defaultVal int) int {
	if val, err := strconv.Atoi(os.Getenv(key)); err == nil && val > 0 {
		return val
	}
	return defaultVal
}
func GetEnvBool(key string, defaultVal bool) bool {
	if val, err := strconv.ParseBool(os.Getenv(key)); err == nil {
		return val
	}
	return defaultVal
}

// GetVideoDuration returns video duration in seconds
func GetVideoDuration(filePath string) (float64, error) {
	cmd := exec.Command(
		"ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		filePath,
	)

	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	durationStr := strings.TrimSpace(string(output))
	return strconv.ParseFloat(durationStr, 64)
}
func GetVideoResolution(filePath string) (int, int, int, error) {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "stream=width,height",
		"-of", "csv=p=0:s=x",
		filePath,
	)

	out, err := cmd.Output()
	if err != nil {
		return 0, 0, 0, err
	}

	var width, height int
	_, err = fmt.Sscanf(string(out), "%dx%d", &width, &height)
	if err != nil {
		return 0, 0, 0, err
	}

	var resolution int
	switch {
	case height <= 144:
		resolution = 144
		break
	case height <= 480:
		resolution = 480
		break
	case height <= 720:
		resolution = 720
		break
	default:
		resolution = height
		break
	}

	return width, height, resolution, nil
}
func GetImageResolution(filePath string) (int, int, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()

	img, _, err := image.DecodeConfig(f)
	if err != nil {
		return 0, 0, err
	}

	return img.Width, img.Height, nil
}