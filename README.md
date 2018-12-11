# aws-cloudwatch-set-logs-expiry

This is a small maintanence tool for keeping cloudwatch expiery taht by default are set to Never to 14 days

>Most of the time its just forgeting to place the expiery.
>So just keep the dollars in yout pocket

To build

```bash
go get github.com/aws/aws-sdk-go
go build -ldflags "-s -w"
```

And execute the binary artifact `./aws-cloudwatch-set-logs-expiry`

Make sure to have a proper permissions to execute it

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "VisualEditor0",
      "Effect": "Allow",
      "Action": [
        "logs:DescribeLogGroups",
        "logs:PutRetentionPolicy"
      ],
      "Resource": "*"
    }
  ]
}
```