local json = require("cjson")

local _M = {
  defaultHTML = [[
    <!DOCTYPE html>
    <html lang="en">

    <head>
        <meta charset="UTF-8">
        <meta http-equiv="X-UA-Compatible" content="IE=edge">
        <meta name="viewport" content="width=device-width, initial-scale=1">
        <meta http-equiv="refresh" content="10">
        <title>Loading...</title>
        <style type="text/css">
            html,
            body {
                width: 100%;
                height: 100%;
            }

            body {
                background: #f8f8f8;
            }

            .content {
                width: 400px;
                height: 300px;
                margin: 0 auto;
                padding: 80px 0 50px;
                display: block;
                position: absolute;
                top: 50%;
                left: 50%;
                margin: -150px 0 0 -200px; 
            }

            .text {
                margin: 0 20px;
                color: #666;
                font-size: 18px;
                line-height: 30px;
                text-align: center;
            }
            .text2 {
                color: #666;
                font-size: 14px;
                margin: 20px 0;
                text-align: center;
            }
            .text2 a {
                color: #666;
            }

            .link {
                width: 285px;
                height: 50px;
                display: block;
                color: #fff;
                font-size: 20px;
                line-height: 50px;
                text-align: center;
                margin: 20px auto 0;
                text-decoration: none;
                border-radius: 4px;
                cursor: pointer;
                background: #1890ff;
            }

            /**/
            .loader {
                width: 140px;
                margin: 50px auto 0;
            }

            .loader-inner {
                margin-left: 40px;
            }

            @-webkit-keyframes rotate_pacman_half_up {
                0% {
                    -webkit-transform: rotate(270deg);
                    transform: rotate(270deg);
                }

                50% {
                    -webkit-transform: rotate(360deg);
                    transform: rotate(360deg);
                }

                100% {
                    -webkit-transform: rotate(270deg);
                    transform: rotate(270deg);
                }
            }

            @keyframes rotate_pacman_half_up {
                0% {
                    -webkit-transform: rotate(270deg);
                    transform: rotate(270deg);
                }

                50% {
                    -webkit-transform: rotate(360deg);
                    transform: rotate(360deg);
                }

                100% {
                    -webkit-transform: rotate(270deg);
                    transform: rotate(270deg);
                }
            }

            @-webkit-keyframes rotate_pacman_half_down {
                0% {
                    -webkit-transform: rotate(90deg);
                    transform: rotate(90deg);
                }

                50% {
                    -webkit-transform: rotate(0deg);
                    transform: rotate(0deg);
                }

                100% {
                    -webkit-transform: rotate(90deg);
                    transform: rotate(90deg);
                }
            }

            @keyframes rotate_pacman_half_down {
                0% {
                    -webkit-transform: rotate(90deg);
                    transform: rotate(90deg);
                }

                50% {
                    -webkit-transform: rotate(0deg);
                    transform: rotate(0deg);
                }

                100% {
                    -webkit-transform: rotate(90deg);
                    transform: rotate(90deg);
                }
            }

            @-webkit-keyframes pacman-balls {
                75% {
                    opacity: 0.7;
                }

                100% {
                    -webkit-transform: translate(-100px, -6.25px);
                    transform: translate(-100px, -6.25px);
                }
            }

            @keyframes pacman-balls {
                75% {
                    opacity: 0.7;
                }

                100% {
                    -webkit-transform: translate(-100px, -6.25px);
                    transform: translate(-100px, -6.25px);
                }
            }

            .pacman {
                position: relative;
            }

            .pacman>div:nth-child(2) {
                -webkit-animation: pacman-balls 1s -0.99s infinite linear;
                animation: pacman-balls 1s -0.99s infinite linear;
            }

            .pacman>div:nth-child(3) {
                -webkit-animation: pacman-balls 1s -0.66s infinite linear;
                animation: pacman-balls 1s -0.66s infinite linear;
            }

            .pacman>div:nth-child(4) {
                -webkit-animation: pacman-balls 1s -0.33s infinite linear;
                animation: pacman-balls 1s -0.33s infinite linear;
            }

            .pacman>div:nth-child(5) {
                -webkit-animation: pacman-balls 1s 0s infinite linear;
                animation: pacman-balls 1s 0s infinite linear;
            }

            .pacman>div:first-of-type {
                width: 0px;
                height: 0px;
                border-right: 20px solid transparent;
                border-top: 20px solid #d4d4d4;
                border-left: 20px solid #d4d4d4;
                border-bottom: 20px solid #d4d4d4;
                border-radius: 20px;
                -webkit-animation: rotate_pacman_half_up 0.5s 0s infinite;
                animation: rotate_pacman_half_up 0.5s 0s infinite;
                position: relative;
                left: -30px;
            }

            .pacman>div:nth-child(2) {
                width: 0px;
                height: 0px;
                border-right: 20px solid transparent;
                border-top: 20px solid #d4d4d4;
                border-left: 20px solid #d4d4d4;
                border-bottom: 20px solid #d4d4d4;
                border-radius: 20px;
                -webkit-animation: rotate_pacman_half_down 0.5s 0s infinite;
                animation: rotate_pacman_half_down 0.5s 0s infinite;
                margin-top: -40px;
                position: relative;
                left: -30px;
            }

            .pacman>div:nth-child(3),
            .pacman>div:nth-child(4),
            .pacman>div:nth-child(5),
            .pacman>div:nth-child(6) {
                background-color: #d4d4d4;
                border-radius: 100%;
                margin: 2px;
                width: 8px;
                height: 8px;
                position: absolute;
                -webkit-transform: translate(0, -6.25px);
                -ms-transform: translate(0, -6.25px);
                transform: translate(0, -6.25px);
                top: 20px;
                left: 70px;
            }
        </style>
        <script type="text/javascript">
            function myrefresh() {
                window.location.reload();
            }
        </script>
    </head>

    <body>
        <div class="content">
            <p class="text">您的应用正在准备中，请稍等一会儿</p>
            <a class="link" onclick="javascript:myrefresh();" style="cursor:pointer;">刷 新</a>
            <div class="loader">
                <div class="loader-inner pacman">
                    <div></div>
                    <div></div>
                    <div></div>
                    <div></div>
                    <div></div>
                </div>
            </div>
            POWER
        </div>
    </body>

    </html>
  ]]
}

function _M.call()
  ngx.header["Content-type"] = "text/html"
  local html = ""
  if (os.getenv("DISABLE_POWER") == "true")
  then
    html = string.gsub(_M.defaultHTML, "POWER", "", 1)
  else
    html = string.gsub(_M.defaultHTML, "POWER", [[<p class="text2" id="power">Power By <a href="https://www.rainbond.com" target="_blank" rel="noopener noreferrer">Rainbond</a></p>]],1)
  end
  ngx.print(html)
  ngx.status = ngx.HTTP_OK
end

return _M
