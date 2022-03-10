<?
header("Access-Control-Allow-Origin: *");
header("Access-Control-Allow-Headers: *");
$dir    = 'data';
$files = scandir($dir);
$result=[];
foreach($files as $file) {
    if ($file!="." && $file!=".."){
        $json = json_decode(file_get_contents($dir."/".$file));
        $result[]=$json;
    }
};

header('Content-Type: application/json');
echo json_encode($result);