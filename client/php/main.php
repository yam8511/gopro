<?php

include 'jsonrpc.php';

$client = new JsonRPC("127.0.0.1", 50052);
$ex = $client->Dial();
if (is_null($ex)) {
    $r = $client->Call("arith.Sum", array('A'=>7, 'B'=>8));
    var_dump($r);
    $r = $client->Call("arith.Sum", array('A'=>7, 'B'=>8));
    var_dump($r);
} else {
    echo "Can't not connect! " . $ex->getMessage();
}
