<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no" />
        <link rel="apple-touch-icon" sizes="76x76" href="/public/img/apple-icon.png" />
        <link rel="icon" type="image/png" href="/public/img/favicon.png" />
        <title>{{ .Title }} &minus; ParrotDisco.pl</title>
        <link href="https://fonts.googleapis.com/css?family=Open+Sans:300,400,600,700" rel="stylesheet" />
        <link href="/public/css/nucleo-icons.css" rel="stylesheet" />
        <link href="/public/css/nucleo-svg.css" rel="stylesheet" />
        <script src="https://kit.fontawesome.com/42d5adcbca.js" crossorigin="anonymous"></script>
        <link href="/public/css/nucleo-svg.css" rel="stylesheet" />
        <link id="pagestyle" href="/public/css/soft-ui-dashboard.css?v=1.0.7" rel="stylesheet" />

        {{template "head" .}}
    </head>

    <body class="g-sidenav-show bg-gray-100">
        {{template "content" .}}

        <script src="/public/js/core/popper.min.js"></script>
        <script src="/public/js/core/bootstrap.min.js"></script>
        <script src="/public/js/plugins/perfect-scrollbar.min.js"></script>
        <script src="/public/js/plugins/smooth-scrollbar.min.js"></script>
        <script src="/public/js/plugins/chartjs.min.js"></script>
        <script>
            var win = navigator.platform.indexOf('Win') > -1;
            if (win && document.querySelector('#sidenav-scrollbar')) {
                var options = {
                    damping: '0.5',
                };
                Scrollbar.init(document.querySelector('#sidenav-scrollbar'), options);
            }
        </script>

        <script async defer src="https://buttons.github.io/buttons.js"></script>
        <script src="/public/js/soft-ui-dashboard.min.js?v=1.0.7"></script>

        {{template "footer" .}}
    </body>
</html>
