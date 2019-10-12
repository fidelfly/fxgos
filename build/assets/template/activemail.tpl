<!DOCTYPE html>

<html>
<head>
    <title>Activate your account</title>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
</head>
<body>
    <div class="container">
        <div class="row">
            <div class="hero-text">
                <h3 class="title">Dear {{.name}},</h3>
                <p class="content">
                    Thank you for using fxgos.
                    <br />
                    <br />
                    Please click the button below to activate your account
                </p>
                <a href="{{.activelink}}"><button type="button">Activate</button></a>
                <p class="content">
                    Thank you very much for your kind attention.
                    <br />
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