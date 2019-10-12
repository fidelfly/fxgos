<!DOCTYPE html>
<html>
<head>
    <title>Reset password</title>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
</head>
<body>
    <div class="container">
        <div class="row">
            <div class="hero-text">
                <h3 class="title">Dear {{.name}},</h3>
                <p class="content">
                    You can use the flowing link to reset your password:
                </p>
                <a href="{{.resetPwdLink}}">Reset password</a>
                <p class="content">
                    If you don't use this link within 24 hours, it will expire.
                </p>
                <p class="content">
                    Thank you very much for your kind attention.
                </p>
                <p class="signature">
                    Best Regards,
                    <br />
                    Fidel Xu
                    <br />
                </p>
            </div>
        </div>
    </div>
</body>
</html>