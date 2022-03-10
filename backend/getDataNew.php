<?
header("Access-Control-Allow-Origin: *");
header("Access-Control-Allow-Headers: *");
$dir    = 'ping';
$files = scandir($dir);
$result=[];
foreach($files as $file) {
    if ($file!="." && $file!=".."){
        $info = array(
            "name" => explode(".",$file)[0],
        );
        $info["ping"] = explode("\n", file_get_contents($dir."/".$file));
        $speed = getSpeedInfo($file);
        if ($speed!=null) {
            $info["download"] = $speed->Download;
            $info["upload"] = $speed->Upload;
        }
        $result[]=$info;
    }
};

function getSpeedInfo($file) {
    if (file_exists("speed/".$file)) {
        $json = json_decode(file_get_contents("speed/".$file));
        return $json;
    }
    return null;
}

header('Content-Type: application/json');
echo json_encode($result);