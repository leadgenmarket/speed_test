<?
if (file_exists('data/')) {
    foreach (glob('data/*') as $file) {
        unlink($file);
    }
}

if (file_exists('speed/')) {
    foreach (glob('speed/*') as $file) {
        unlink($file);
    }
}

if (file_exists('ping/')) {
    foreach (glob('ping/*') as $file) {
        unlink($file);
    }
}