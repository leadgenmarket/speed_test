<?
header("Access-Control-Allow-Origin: *");
header("Access-Control-Allow-Headers: *");
if ($_SERVER['REQUEST_METHOD'] === 'GET'){
    $settings = json_decode(file_get_contents("settings.json"));
    header('Content-Type: application/json; charset=utf-8');
    echo json_encode($settings);
} else if ($_SERVER['REQUEST_METHOD'] === 'POST') {
    include "auth.php";
    if (checkAuthAction($key)) {
        $postData = file_get_contents('php://input');
        file_put_contents("settings.json",$postData);
        header('Content-Type: application/json; charset=utf-8');
        echo json_encode(array("result"=>"ok"));
    } else {
        header('Content-Type: application/json; charset=utf-8');
        echo json_encode(array("result"=>"not authorized"));
    }
}
