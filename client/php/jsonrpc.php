<?php
class JsonRPC
{
    private $host;
    private $port;
    private $conn;
    private $timeout;

    public function __construct($host, $port, $timeout = 5)
    {
        $this->host = $host;
        $this->port = $port;
        $this->conn = NULL;
        $this->timeout = $timeout;
    }

    public function Dial()
    {
        $this->conn = fsockopen($this->host, $this->port, $errno, $errstr, 3);
        if (!$this->conn) {
            $this->conn = NULL;
            return new \Exception("JsonRPC Dial Failed: $errstr ($errno)");
        }
        return NULL;
    }

    public function Call($method, $params, $id = 1)
    {
        if (!$this->conn) {
            $ex = $this->Dial();
            if (!is_null($ex)) {
                return $ex;
            }
        }
        $err = fwrite($this->conn, json_encode(array(
                'method' => $method,
                'params' => array($params),
                'id'     => $id,
            ))."\n");
        if ($err === false) {
            return new \Exception("JsonRPC Send Failed");
        }
        stream_set_timeout($this->conn, $this->timeout);
        for (;;) {
            // 檢查回傳            
            $line = fgets($this->conn);
            if ($line !== false) {
                $this->Close();
                return json_decode($line, true);
            }

            // 檢查Timeout
            $info = stream_get_meta_data($this->conn);
            if ($info['timed_out']) {
                $this->Close();
                return new \Exception("JsonRPC Init Time Out");
            }
		}
    }

    public function SetTimeout($t = 5)
    {
        $this->timeout = $t;
        return $this;
    }

    public function Close() {
        fclose($this->conn);
        $this->conn = NULL;
    }
}