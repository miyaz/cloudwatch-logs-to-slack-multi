
1. update version

modify template.yml

```diff:template.yml
-    SemanticVersion: 0.0.1
+    SemanticVersion: 0.0.2
     # best practice is to use git tags for each release and link to the version tag as your source code URL
-    SourceCodeUrl: https://github.com/miyaz/cloudwatch-logs-to-slack-multi/tree/0.0.1
+    SourceCodeUrl: https://github.com/miyaz/cloudwatch-logs-to-slack-multi/tree/0.0.2
```

2. git push

```
git add ~
git commit ~
git push origin master
git tag 0.0.2
git push origin 0.0.2
```

3. build/package/publish

```
sam build

sam package \
    --profile {profile} --region ap-northeast-1 \
    --template-file .aws-sam/build/template.yaml \
    --output-template-file packaged.yaml \
    --s3-bucket miyaz-package-for-serverlessrepo

sam publish \
    --profile {profile} --region ap-northeast-1 \
    --template packaged.yaml
```

need this bucket policy:

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": {
                "Service":  "serverlessrepo.amazonaws.com"
            },
            "Action": "s3:GetObject",
            "Resource": "arn:aws:s3:::miyaz-package-for-serverlessrepo/*"
        }
    ]
}
```

