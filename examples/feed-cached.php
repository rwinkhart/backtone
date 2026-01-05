<?php
$command = "/usr/local/bin/backtone";
$argument = "BASE64ENCODEDINPUTXML";
$cacheFile = __DIR__ . "/cache/feed.xml";
$tempFile = "/tmp/feed.xml" . ".tmp";

// serve cached file immediately
header("Content-Type: application/xml; charset=utf-8");
readfile($cacheFile);

// trigger background regeneration
exec(escapeshellcmd($command) . " " . escapeshellarg($argument) . " > " . escapeshellarg($tempFile) . " 2>&1 && mv " . escapeshellarg($tempFile) . " " . escapeshellarg($cacheFile) . " &");
?>