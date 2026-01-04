<?php
$command = "/usr/local/bin/backtone";
$argument = "BASE64ENCODEDINPUTXML";

$output = shell_exec(escapeshellcmd($command) . " " . escapeshellarg($argument));
header("Content-Type: application/xml; charset=utf-8");
echo $output;
?>
