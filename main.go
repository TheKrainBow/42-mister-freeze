package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"mister-freeze/apiManager"
	"mister-freeze/config"
)

const CONFIG_PATH = "./config.yml"

func init42FreezeAPI() error {
	APIClient, err := apiManager.NewAPIClient(config.FTFreeze, apiManager.APIClientInput{
		AuthType:     apiManager.AuthTypePassword,
		TokenURL:     config.ConfigData.Freeze42.TokenUrl,
		Endpoint:     config.ConfigData.Freeze42.Endpoint,
		TestPath:     config.ConfigData.Freeze42.TestPath,
		ClientID:     config.ConfigData.Freeze42.Uid,
		ClientSecret: config.ConfigData.Freeze42.Secret,
		Username:     config.ConfigData.Freeze42.Username,
		Password:     config.ConfigData.Freeze42.Password,
	})
	if err != nil {
		return fmt.Errorf("couldn't create freeze api client: %w", err)
	}
	err = APIClient.TestConnection()
	if err != nil {
		return fmt.Errorf("api connection test to freeze failed: %w", err)
	}
	return nil
}

func init42v2API() error {
	APIClient, err := apiManager.NewAPIClient(config.FTv2, apiManager.APIClientInput{
		AuthType:     apiManager.AuthTypeClientCredentials,
		TokenURL:     config.ConfigData.ApiV2.TokenUrl,
		Endpoint:     config.ConfigData.ApiV2.Endpoint,
		TestPath:     config.ConfigData.ApiV2.TestPath,
		ClientID:     config.ConfigData.ApiV2.Uid,
		ClientSecret: config.ConfigData.ApiV2.Secret,
		Scope:        config.ConfigData.ApiV2.Scope,
	})
	if err != nil {
		return fmt.Errorf("couldn't create 42v2 api client: %w", err)
	}
	err = APIClient.TestConnection()
	if err != nil {
		return fmt.Errorf("api connection test to 42v2 failed: %w", err)
	}
	return nil
}

type RequestData struct {
	UserIDs            []string `json:"user_ids"`
	UserExcludedIDs    []string `json:"excluded_user_ids,omitempty"`
	BeginDate          string   `json:"begin_date"`
	ExpectedEndDate    string   `json:"expected_end_date"`
	Reason             string   `json:"reason"`
	IsFreeFreeze       bool     `json:"is_free_freeze"`
	StudentDescription string   `json:"student_description"`
	StaffDescription   string   `json:"staff_description"`
}

type User struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	IsActive bool   `json:"active?"`
	IsStaff  bool   `json:"staff?"`
}

type QuestUser struct {
	User User `json:"user"`
}

func getAllUsers(data RequestData) RequestData {
	cleanDuplicate := map[string]bool{}
	page := 0
	for {
		url := fmt.Sprintf("/quests/37/quests_users?filter[campus_id]=41&filter[validated]=false&page[size]=100&page[number]=%d", page)
		resp, err := apiManager.GetClient(config.FTv2).Get(url)
		if err != nil {
			os.Exit(1)
		}
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			os.Exit(1)
		}
		var res []QuestUser
		err = json.Unmarshal(bodyBytes, &res)
		if err != nil {
			log.Printf("%s\n", err.Error())
			os.Exit(1)
		}

		for _, qu := range res {
			skip := false
			for _, s := range data.UserExcludedIDs {
				parsedID, _ := strconv.Atoi(s)
				if qu.User.ID == parsedID || qu.User.Login == s || qu.User.IsStaff {
					fmt.Printf("%s is excluded\n", qu.User.Login)
					skip = true
					break
				}
			}

			if skip {
				continue
			}

			if qu.User.IsActive {
				if _, exists := cleanDuplicate[qu.User.Login]; !exists {
					data.UserIDs = append(data.UserIDs, qu.User.Login)
					cleanDuplicate[qu.User.Login] = true
				}
			}
		}
		if len(res) < 100 {
			break
		}
		fmt.Printf("Found %d users on page %d\n", len(res), page)
		page++
	}
	fmt.Printf("Found %d users\n", len(data.UserIDs))
	return data
}

func main() {
	config.LoadConfig(CONFIG_PATH)
	err := init42v2API()
	if err != nil {
		log.Fatalf("couldn't start 42 v2 API: %s\n", err.Error())
	}
	err = init42FreezeAPI()
	if err != nil {
		log.Fatalf("couldn't start 42 Freeze API: %s\n", err.Error())
	}

	data := askUserPayloadInfo()
	data = getAllUsers(data)

	fmt.Println("\nCollected data:")
	printData(data)

	if askForValidBool(fmt.Sprintf("This freeze will hit %d users, continue? [y/N] ", len(data.UserIDs))) {
		fmt.Printf("Validated\n")
		data.UserExcludedIDs = []string{}
		resp, err := apiManager.GetClient(config.FTFreeze).Post("/freezes/compensation/bulk", data)
		if err != nil {
			log.Fatalf("Couldn't Post Freeze: %s", err.Error())
		} else {
			log.Printf("Response: %s\n", resp.Status)
		}
	} else {
		fmt.Printf("Operation canceled\n")
	}
}

func askUserPayloadInfo() RequestData {
	var data RequestData

	data.BeginDate = askForValidDate("Enter begin_date (YYYY-MM-DD): ")
	data.ExpectedEndDate = askForValidEndDate("Enter expected_end_date (YYYY-MM-DD): ", data.BeginDate)
	data.Reason = askForValidReason("Enter reason (other, personnal, professional, medical): ")
	data.IsFreeFreeze = askForValidBool("Is this a free freeze [y/N]: ")
	data.StudentDescription = askForNonEmptyField("Enter student description: ")
	data.StaffDescription = askForNonEmptyField("Enter staff description: ")
	data.UserExcludedIDs = askForExcludedUserList("Enter logins you want to EXCLUDE: ")
	return (data)
}

func askForExcludedUserList(prompt string) []string {
	var input string
	scanner := bufio.NewScanner(os.Stdin)
	for {
		// Print prompt
		fmt.Print(prompt)
		scanner.Scan()
		input = scanner.Text()

		userList := splitAndTrim(input)
		return userList
	}
}

// splitAndTrim splits the input by commas and trims spaces from each element
func splitAndTrim(input string) []string {
	// Split by commas and trim each element
	var result []string
	for _, user := range strings.Split(input, ",") {
		// Trim leading and trailing spaces
		trimmedUser := strings.TrimSpace(user)
		if trimmedUser != "" {
			result = append(result, trimmedUser)
		}
	}
	return result
}

func askForValidReason(prompt string) string {
	validReasons := []string{"other", "personnal", "professional", "medical"}
	var input string
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(prompt)
		scanner.Scan()
		input = scanner.Text()
		input = strings.ToLower(input)
		if contains(validReasons, input) {
			return input
		} else {
			fmt.Println("Invalid reason. Please choose from: other, personnal, professional, medical.")
		}
	}
}

func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

// askForValidDate asks the user to input a date and ensures it's in YYYY-MM-DD format
func askForValidDate(prompt string) string {
	var input string
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(prompt)
		scanner.Scan()
		input = scanner.Text()
		if isValidDate(input) {
			return input
		} else {
			fmt.Println("Invalid date format. Please enter a date in YYYY-MM-DD format.")
		}
	}
}

func askForValidEndDate(prompt string, beginDate string) string {
	var input string
	beginDateTime, _ := time.Parse("2006-01-02", beginDate)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(prompt)
		scanner.Scan()
		input = scanner.Text()
		if isValidDate(input) {
			endDateTime, _ := time.Parse("2006-01-02", input)
			if endDateTime.After(beginDateTime) || endDateTime.Equal(beginDateTime) {
				return input
			} else {
				fmt.Println("Expected end date cannot be before begin date. Please enter a valid date.")
			}
		} else {
			fmt.Println("Invalid date format. Please enter a date in YYYY-MM-DD format.")
		}
	}
}

// isValidDate checks if a date string is in the format YYYY-MM-DD
func isValidDate(date string) bool {
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}

// askForValidBool asks for a boolean value and validates input
func askForValidBool(prompt string) bool {
	var input string
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(prompt)
		scanner.Scan()
		input = scanner.Text()
		input = strings.ToLower(input)
		if isValidBool(input) {
			return parseBool(input)
		} else {
			fmt.Println("Invalid input. Please enter true/false, t/f, yes/no, or their abbreviations.")
		}
	}
}

// isValidBool checks if the input is a valid boolean representation
func isValidBool(input string) bool {
	validBools := []string{"true", "false", "t", "f", "yes", "no", "y", "n"}
	for _, v := range validBools {
		if input == v {
			return true
		}
	}
	return false
}

// parseBool converts user input into a boolean value
func parseBool(input string) bool {
	switch input {
	case "true", "t", "yes", "y":
		return true
	case "false", "f", "no", "n":
		return false
	default:
		return false
	}
}

// askForNonEmptyField ensures the user input is not empty
func askForNonEmptyField(prompt string) string {
	var input string
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(prompt)
		scanner.Scan()
		input = scanner.Text()
		if input != "" {
			return input
		} else {
			fmt.Println("This field cannot be empty. Please enter a value.")
		}
	}
}

// printData prints the collected data
func printData(data RequestData) {
	dataJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatal("Error marshaling data: ", err)
	}

	// Print the pretty-printed JSON (like jq output)
	fmt.Println(string(dataJSON))
}
