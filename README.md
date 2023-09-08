# Copy Image A2B

## Configuration

Before running the program, you need to configure it to connect to the MongoDB databases and specify HTTP download settings. You can do this by creating a .env configuration file with the following example configuration:

```env
A_HOSTNAME=https://DEV_HOSTNAME/
A_MONGO_HOST=DEV_MONGO_HOST
B_MONGO_HOST=PROD_MONGO_HOST
```

Please modify the above information according to your environment.

## Usage
Once you have configured the program, follow these steps to run it:

1. In the program's root directory, use the following command to run the program:
    ```sh
    go run main.go
    ```

2. The program will connect to Prod MongoDB, check if the image path exists locally. If the image is not present locally, it will attempt to download the image from Dev environment via HTTP.

3. Downloaded images will be stored in the local folder specified by download_path.

## Contribution
If you find any issues or have improvement suggestions, feel free to submit an issue or pull request to this repository.

## License
This program is released under the MIT License. Detailed information can be found in the [LICENSE](./license) file.