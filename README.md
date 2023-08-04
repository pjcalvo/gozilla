# Gozilla - Simple Load Testing Library for Go

Gozilla is a lightweight and user-friendly load testing library for Go, designed to help you simulate and evaluate the performance of your web applications. With Gozilla, you can easily run load tests with customizable configurations and define expected behavior for each request.

## Installation

To use Gozilla in your Go project, simply import it:

```bash
go get -u github.com/pjcalvo/gozilla/pkg
```

## Getting started

To start using Gozilla, you can follow these simple steps:

1. Import the Gozilla package into your code:

```go
import (
    "errors"
    "log"
    "net/http"
    "time"

    gozilla "github.com/pjcalvo/gozilla/pkg"
)
```

2. Create and configure a test suite with the desired load testing parameters:

```go
suite := gozilla.NewTestSuite().
    WithBaseURL("https://pjcalvo.github.io").
    WithDuration(time.Minute * 5).
    WithUsers(3)
    WithThinkTime(time.Second * 1)
```

3. Define the test tasks you want to execute during the load test:

```go
reqHomePage, _ := http.NewRequest("GET", "/", nil)
reqAbout, _ := http.NewRequest("GET", "/about", nil)

tasks := []gozilla.Task{
    {
        Request: reqHomePage,
        Label:   "homepage",
    },
    {
        Request: reqAbout,
        Label:   "about",
        ExpectedFunc: func(r *http.Response) error {
            if r.StatusCode != 200 {
                return errors.New("response code is not expected")
            }
            return nil
        },
    },  
}
```

4. Run the load test by executing the test suite and passing the defined tasks:

```go
if err := suite.ExecuteTest(tasks); err != nil {
    log.Fatal(err)
}
```

## Results Reporting

Upon executing the load test, Gozilla will provide you with insightful results in the form of a timestamped CSV report file. The report (`<timestamp>_results.csv`) will contain detailed information about each request's response time, status code, and any encountered errors during the test. This allows you to analyze the performance of your web application comprehensively.

Additionally, Gozilla features a built-in plotter that will be displayed on the console. The plotter provides a visual representation of the load test's performance, giving you an intuitive understanding of how your application behaves under simulated user traffic. This graphical output facilitates quick identification of performance bottlenecks and helps you optimize your application for better scalability and user experience.


```csv
timestamp,threadID,status,duration (ms),label,responseMessage,bytes,success,error
2023-08-04T14:42:43+02:00,0,200,771.00,homepage,,29933,true,
2023-08-04T14:42:44+02:00,0,200,361.00,about,,20451,true,
```

![Plotter](https://github.com/pjcalvo/gozilla/raw/main/resources/plotter.png)

## Contributing

We welcome contributions to Gozilla! If you find any issues or have suggestions for improvements, please feel free to open an issue or submit a pull request on the GitHub repository: [https://github.com/pjcalvo/gozilla](https://github.com/pjcalvo/gozilla)

## License

Gozilla is open-source software licensed under the MIT License. See the LICENSE file for more details.

## Acknowledgments

Gozilla was inspired by the need for a simple and efficient load testing library in Go. Also inspired in a sense by [locustIO](https://locust.io/).

