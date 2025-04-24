````markdown
# MediBridge API

MediBridge is a healthcare management system designed to allow users to register, manage patient data, diagnose, and handle allergies, conditions, and vitals. This API provides endpoints for user authentication, patient management, and more.

## Table of Contents

- [API Overview](#api-overview)
- [Features](#features)
- [Setup and Installation](#setup-and-installation)
- [Running the Application](#running-the-application)
- [API Endpoints](#api-endpoints)
  - [User Authentication](#user-authentication)
  - [Patient Management](#patient-management)
  - [Diagnoses Management](#diagnoses-management)
  - [Conditions Management](#conditions-management)
  - [Allergy Management](#allergy-management)
  - [Vitals Management](#vitals-management)
- [Swagger Documentation](#swagger-documentation)
- [Error Handling](#error-handling)
- [Contributing](#contributing)

## API Overview

MediBridge API provides functionalities for managing patient information, user authentication, and health records in the healthcare domain. It uses the Go programming language with the Chi router for routing and middleware. The API also integrates Swagger for interactive API documentation.

## Features

- **User Authentication**: Signup and signin routes to authenticate users.
- **Patient Management**: Register, list, and manage patient data.
- **Diagnoses Management**: Add, update, and delete diagnoses related to patients.
- **Conditions Management**: Handle conditions associated with patients.
- **Allergy Management**: Record and update patient allergies.
- **Vitals Management**: Capture, update, and delete patient vitals.
- **Swagger Documentation**: Interactive API documentation with Swagger UI.

## Setup and Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/your-username/medibridge.git
   cd medibridge
   ```
````

2.  **Install dependencies**:

    Ensure Go and required dependencies are installed. You can install Go from [here](https://golang.org/dl/).

    Run the following commands to install dependencies:

    ```bash
    go mod tidy

    ```

3.  **Configure Environment Variables**:

    Make sure to set up your environment variables for any necessary configuration (e.g., database URL, JWT secrets). You can use `.env` or any preferred method for environment variables.

    Example `.env` file:

    ```env
    DB_URL=your_database_url
    JWT_SECRET=your_jwt_secret

    ```

## Running the Application

To run the API locally, you can use the following command:

```bash
go run ./cmd/main.go

```

This will start the API server on `http://localhost:8080`.

### Running with Docker

If you prefer using Docker, build and run the application with the following commands:

```bash
docker build -t medibridge .
docker run -p 8080:8080 medibridge

```

## API Endpoints

### User Authentication

#### POST /v1/signup

Register a new user.

**Request body**:

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response**:

- `201 Created` on success.
- `400 Bad Request` if input is invalid.

#### POST /v1/signin

Sign in to obtain a JWT token.

**Request body**:

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response**:

- `200 OK` with JWT token on success.
- `401 Unauthorized` if credentials are incorrect.

---

### Patient Management

#### POST /v1/patient

Register a new patient.

**Request body**:

```json
{
  "name": "John Doe",
  "dob": "1980-01-01",
  "gender": "male"
}
```

**Response**:

- `201 Created` on success.
- `400 Bad Request` if input is invalid.

#### GET /v1/patient/list

Get a list of all registered patients.

**Response**:

- `200 OK` with a list of patients.

---

### Diagnoses Management

#### POST /v1/patient/{patientID}/diagnoses

Add a diagnosis for a patient.

**Request body**:

```json
{
  "diagnosis": "Hypertension",
  "details": "Patient has high blood pressure."
}
```

**Response**:

- `201 Created` on success.
- `400 Bad Request` if input is invalid.

---

### Conditions Management

#### POST /v1/patient/{patientID}/condition

Add a medical condition for a patient.

**Request body**:

```json
{
  "condition": "Asthma",
  "severity": "moderate"
}
```

**Response**:

- `201 Created` on success.
- `400 Bad Request` if input is invalid.

---

### Allergy Management

#### POST /v1/patient/{patientID}/allergy

Record an allergy for a patient.

**Request body**:

```json
{
  "allergy": "Peanut",
  "reaction": "Anaphylaxis"
}
```

**Response**:

- `201 Created` on success.
- `400 Bad Request` if input is invalid.

---

### Vitals Management

#### POST /v1/patient/{patientID}/vitals

Capture vitals for a patient.

**Request body**:

```json
{
  "temperature": 98.6,
  "blood_pressure": "120/80",
  "heart_rate": 75
}
```

**Response**:

- `201 Created` on success.
- `400 Bad Request` if input is invalid.

---

## Swagger Documentation

The Swagger documentation for the API is available at:

```
http://localhost:8080/swagger/index.html

```

This will give you a detailed, interactive UI where you can explore and test the API endpoints directly.

---

## Error Handling

The API uses standard HTTP status codes for error handling:

- `200 OK`: Request was successful.
- `201 Created`: Resource was created successfully.
- `400 Bad Request`: Invalid input or missing required fields.
- `401 Unauthorized`: Authentication failed or user lacks necessary permissions.
- `404 Not Found`: The requested resource could not be found.
- `500 Internal Server Error`: An unexpected error occurred.

## Contributing

We welcome contributions to improve MediBridge! If you'd like to contribute, please fork the repository and create a pull request.

1.  Fork the repository.
2.  Create a new branch (`git checkout -b feature-name`).
3.  Commit your changes (`git commit -am 'Add new feature'`).
4.  Push to the branch (`git push origin feature-name`).
5.  Create a new Pull Request.

---

## License

This project is licensed under the MIT License - see the [LICENSE](https://chatgpt.com/c/LICENSE) file for details.

```

### Key Points Covered:
- **API Overview**: Introduction to what the API does.
- **Features**: Lists key features like user authentication, patient management, and health data handling.
- **Setup and Installation**: Instructions to set up the project locally or using Docker.
- **Running the Application**: Explains how to start the server both manually and with Docker.
- **API Endpoints**: Provides a detailed breakdown of the various routes available in the API, including request bodies and expected responses.
- **Swagger Documentation**: Includes a link to access the API docs.
- **Error Handling**: Outlines the HTTP status codes used in the API.
- **Contributing**: Explains how to contribute to the project.

Let me know if you need any further adjustments!

```
