<?php
file_put_contents("test", json_encode($_REQUEST));
if($_REQUEST["key"]!="swGh889KyxjWyz") die();
unset($_REQUEST["key"]);
if ($_REQUEST["type"]=="ping") {
    file_put_contents("ping/".$_REQUEST["user"].".json", $_REQUEST["msg"]."\n", FILE_APPEND);
} else if ($_REQUEST["type"]=="speed") {
    file_put_contents("speed/".$_REQUEST["user"].".json", $_REQUEST["msg"]);
}