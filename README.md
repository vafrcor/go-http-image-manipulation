# go-http-image-manipulation
HTTP Implementation of Image Manipulation using Golang and OpenCV


## Basic Instructions
- Installation
    - Install latest Golang version (e.g. `go1.21.x`)
    - Install [OpenCV](https://gocv.io/getting-started/) => `>= 4.7.x`
    - Install application dependencies 
        - Run `go mod download` to download dependencies
        - Run `go mod vendor` to set dependencies into project directory
- Build Project
    - Run `go build .` (the binary output => `go-http-image-manipulation`)
- Run HTTP Service
    - From Source => `go run server.go` 
    - From Binary => `./go-http-image-manipulation`
- Unit Tests
    - Run `go test`
    - Run `go test -cover` (show coverage)

## Endpoint Tests 
### Convert image files from PNG to JPEG 
- URL: `[POST] http://localhost:9000/image-png-to-jpeg` 
- Request 
    - Content Type: `multipart/form-data`
    - Fields:

    | Name  | Mandatory  |  Description |
    |:---|:---:|:---|
    | file | yes | image file (`image/png`) |

- Response
    - Content Type: `application/json`
    - Fields:

    | Name  | Type  |  Description |
    |:---|:---:|:---|
    | message | string | detailed message (for both success and error) |
    | status | boolean | `true | false` |  
    | output | string | output path (for preview) | 

    - Examples:
        - Success
        ```json
        {
            "message": "Ok",
            "status": true,
            "output": "http://localhost:9000/static/small-1710681145040310000-100.jpeg"
        }
        ```
        - Error

        ```json
        {
            "message": "only accept image using specific format (png)",
            "status": false,
            "output": ""
        }
        ```

### Resize images according to specified dimensions 
- URL: `[POST] http://localhost:9000/image-resize` 
- Request 
    - Content Type: `multipart/form-data`
    - Fields:

    | Name  | Mandatory  |  Description |
    |:---|:---:|:---|
    | file | yes | image file (`image/png`, `image/jpg`, `image/jpeg`) |
    | width | yes | desired width (`in pixel`) |
    | height | yes | desired height (`in pixel`) |

- Response
    - Content Type: `application/json`
    - Fields:

    | Name  | Type  |  Description |
    |:---|:---:|:---|
    | message | string | detailed message (for both success and error) |
    | status | boolean | `true | false` |  
    | output | string | output path (for preview) | 

    - Example:
        - Success
        ```json
        {
            "message": "Ok",
            "status": true,
            "output": "http://localhost:9000/static/medium-1710685243638707000-100.png"
        }
        ```
        - Error
        ```json
        {
            "message": "only accept image using specific format (png,jpg,jpeg)",
            "status": false,
            "output": ""
        }
        ```

### Compress images to reduce file size while maintaining reasonable quality
- URL: `[POST] http://localhost:9000/image-compression` 
- Request 
    - Content Type: `multipart/form-data`
    - Fields:

    | Name  | Mandatory  |  Description |
    |:---|:---:|:---|
    | file | yes | image file (`image/png`, `image/jpg`, `image/jpeg`) |
    | quality | yes | desired quality (`0 - 100`) |

- Response
    - Content Type: `application/json`
    - Fields:

    | Name  | Type  |  Description |
    |:---|:---:|:---|
    | message | string | detailed message (for both success and error) |
    | status | boolean | `true | false` |  
    | output | string | output path (for preview) |  

    - Example:
        - Success
        ```json
        {
            "message": "Ok",
            "status": true,
            "output": "http://localhost:9000/static/medium-1710686662823893000-70.jpeg"
        }
        ```
        - Error
        ```json
        {
            "message": "http: no such file",
            "status": false,
            "output": ""
        }
        ```

## References
- GoCV
    - [Official](https://gocv.io/)
    - [Go Package](https://pkg.go.dev/gocv.io/x/gocv@v0.35.0)
- [OpenCV](https://docs.opencv.org/4.x/d4/da8/group__imgcodecs.html) 

## Others
- [Sample Input Images](https://drive.google.com/drive/folders/1jnlIXRc6GhXAmOEXwd8I3Iodbuk5T0Hy?usp=sharing)
