<?
if($_REQUEST["key"]!="swGh889KyxjWyz") die();
unset($_REQUEST["key"]);
$data = $_REQUEST;
$data["time"] = time();
file_put_contents("data/".$_REQUEST["user"].".json", json_encode($data));
