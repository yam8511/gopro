package main

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

func deployMaster(node, token string) error {
	if isLocalIP(node) {
		return setupK3sMaster(node, token)
	}

	return nil
}

func deployCluster(nodeMaster, token string, serverNodes, agentNodes []string) error {

	nodes := getIPs()
	errCh := make(chan error)
	for _, nodeIP := range nodes {
		go func(nodeIP string) {
			if inNodes(nodeIP, serverNodes) {
				if isLocalIP(nodeIP) {
					errCh <- setupK3sServer(nodeMaster, nodeIP, token)
				} else {
				}
			} else if inNodes(nodeIP, serverNodes) {
				if isLocalIP(nodeIP) {
					errCh <- setupK3sAgent(nodeMaster, nodeIP, token)
				} else {
				}
			} else {
				errCh <- nil
			}
		}(nodeIP)
	}

	txt := ""
	for i := 0; i < len(nodes); i++ {
		err := <-errCh
		if err != nil {
			txt += err.Error() + "\n"
		}
	}

	if txt != "" {
		return errors.New(txt)
	}

	return nil
}

func setupK3sMaster(ip, token string) error {
	env := []string{
		"server",
		"--cluster-init",
		"--node-ip=" + ip,
		"--token=" + token,
		"--docker",
		"--write-kubeconfig-mode=644",
	}

	cmd := exec.Command("sh", INSTALL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Env = append(cmd.Env,
		"INSTALL_K3S_SKIP_DOWNLOAD=true",
		`INSTALL_K3S_EXEC="`+strings.Join(env, " ")+`"`,
	)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func setupK3sServer(nodeMaster, nodeIP, token string) error {
	env := []string{
		"server",
		"--server=https://" + nodeMaster + ":6443",
		"--node-ip=" + nodeIP,
		"--token=" + token,
		"--docker",
		"--write-kubeconfig-mode=644",
	}

	cmd := exec.Command("sh", INSTALL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Env = append(cmd.Env,
		"INSTALL_K3S_SKIP_DOWNLOAD=true",
		`INSTALL_K3S_EXEC="`+strings.Join(env, " ")+`"`,
	)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func setupK3sAgent(nodeMaster, nodeIP, token string) error {
	env := []string{
		"agent",
		"--server=https://" + nodeMaster + ":6443",
		"--node-ip=" + nodeIP,
		"--token=" + token,
		"--docker",
		"--write-kubeconfig-mode=644",
	}

	cmd := exec.Command("sh", INSTALL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Env = append(cmd.Env,
		"INSTALL_K3S_SKIP_DOWNLOAD=true",
		`INSTALL_K3S_EXEC="`+strings.Join(env, " ")+`"`,
	)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
