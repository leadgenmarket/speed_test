<?
/*ini_set('error_reporting', E_ALL);
ini_set('display_errors', 1);
ini_set('display_startup_errors', 1);*/
include "vendor/autoload.php";


use \Firebase\JWT\JWT;

$key = "H8m3(m{_ZgTp";
$_POST = json_decode(file_get_contents("php://input"),true);
if ($_POST["action"] == "login"){
    authorize($_POST, $key); 
}

if ($_POST["action"] == "check"){
    checkAuth($key); 
}

function authorize($req, $key) {
   
    if($req["login"]=="admin" && $req["pass"]=="Qwerty123"){
        $payload = array(
            "login" => $req["login"],
        );
        $jwt = JWT::encode($payload, $key);
        //file_get_contents(token);
        setcookie("token", $jwt, time()+3600);
        header("Access-Control-Allow-Origin: *");
        header("Access-Control-Allow-Headers: *");
        header('Content-type: application/json');
        http_response_code(200);
        echo json_encode(array("authorized"=>$jwt));
        return;
    }
    header("Access-Control-Allow-Origin: *");
    header("Access-Control-Allow-Headers: *");
    header('Content-type: application/json');
    http_response_code(401);
    echo json_encode(array("authorized"=>false));
}
function checkAuth($key){
    
    header('Content-type: application/json');
    if (isset($_COOKIE['token'])){
        $decoded = JWT::decode($_COOKIE['token'], $key, array('HS256'));
        if ($decoded){
            echo json_encode(array("authorized"=>true));
        } else {
            echo json_encode(array("authorized"=>false));
        }
    }
}

function checkAuthAction($key){
    header("Access-Control-Allow-Origin: *");
    header("Access-Control-Allow-Headers: *");
    header('Content-type: application/json');
    if (isset($_COOKIE['token'])){
        $decoded = JWT::decode($_COOKIE['token'], $key, array('HS256'));
        if ($decoded){
           return true;
        } else {
            unset($_COOKIE['token']);
            return false;
        }
    }
    return false;
}
