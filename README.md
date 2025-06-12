
# 42 Freeze API Tool

This Go application interacts with the **42 Freeze API** to collect user data from the **42 API**, perform operations on users (such as exclusion and freeze), and send the results as a JSON payload to the **42 Freeze API**.

## Prerequisites

To run this project, you'll need:

- Go 1.16 or later
- The `mister-freeze` package dependencies (referenced in the code)
- A valid `config.yml` file containing the necessary API credentials for the **42 API** and **42 Freeze API**

## Setup

### Step 1: Install Dependencies

Ensure you have Go installed. You can download and install Go from the official [Go website](https://golang.org/dl/).

Once you have Go installed, install the required dependencies:

```bash
go mod tidy
```

This will download the necessary dependencies defined in `go.mod`.

### Step 2: Configuration

You will need to configure your API credentials. The configuration file is expected to be at `./config.yml`.

The configuration should contain the following sections:

- **FTFreeze**: Credentials and URL for the **42 Freeze API**
- **FTv2**: Credentials and URL for the **42 API**

Example configuration file (`config.yml`):

```yaml
Freeze42:
  TokenUrl: "https://example.com/token"
  Endpoint: "https://example.com/endpoint"
  TestPath: "/test"
  Uid: "your-uid"
  Secret: "your-secret"
  Username: "your-username"
  Password: "your-password"

ApiV2:
  TokenUrl: "https://api.example.com/token"
  Endpoint: "https://api.example.com/endpoint"
  TestPath: "/test"
  Uid: "your-api-uid"
  Secret: "your-api-secret"
  Scope: "your-api-scope"
```

### Step 3: Build and Run the Application

To build and run the application:

1. Build the Go application:

```bash
go build -o freeze-tool
```

2. Run the application:

```bash
./freeze-tool
```

### Step 4: Interact with the Application

The application will guide you through the following process:

1. **Collecting User Data**: 
   - It will prompt you to enter various details such as:
     - `begin_date` and `expected_end_date`
     - The reason for the freeze (either `other`, `personnal`, `professional`, or `medical`)
     - Whether this is a free freeze (yes/no)
     - Student and staff descriptions

2. **Excluding Users**: 
   - You will be asked to enter user logins or IDs to exclude from the freeze.

3. **Confirming the Freeze**: 
   - The application will display the number of users it will affect, and ask if you want to proceed with the freeze.

4. **Sending the Freeze Request**: 
   - If you confirm, the tool will send a POST request to the **42 Freeze API** to initiate the bulk freeze operation.

### Step 5: Error Handling

If the application encounters an error while interacting with the **42 APIs** (e.g., network failure, invalid credentials), it will log the error and exit with a relevant message.

### Notes

- The **42 Freeze API** and **42 API** must be correctly configured in the `config.yml` file for successful interaction.
- Ensure that you have permission to interact with the **42 APIs** and to perform the actions specified by the tool.
- This tool uses `fmt.Println()` to log the process and `log.Fatalf()` for critical errors.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.