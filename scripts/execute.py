import os, sys, json, time
import subprocess, hashlib
from datetime import datetime

import psutil, requests

####
def rfc3339ms(t: datetime.time):
    return t.astimezone().isoformat(timespec='milliseconds')

def cmd_json(cmd):
    return json.dumps({"commandline": cmd}).replace(" ", "")

def cmd_md5(cmd):
    cmdJSON = cmd_json(cmd)
    hash_object = hashlib.md5(cmdJSON.encode('utf8'))
    return cmdJSON, hash_object.hexdigest()


def reportTask(key, value, detail=""):
   api, taskId = "{{.Api}}", "{{.TaskId}}"

   if api.startswith("{{") or taskId.startswith("{{"): return
   url = "{}?taskId={}&key={}&value={}&detail={}".format(api, key, value, detail)
   res = request(url)

   print("post {}: {}".format(res.status_code, url))


def get_process(cmd):
    _, hv = cmd_md5(cmd)
    pidFile = hv + ".pid"

    if not os.path.exists(pidFile):
        return None, "Pid file not found: {}".format(pidFile)

    return process_from_pidFile(pidFile)


def process_from_pidFile(pidFile):
    with open(pidFile) as f:
        pid = int(f.readline())

    process = psutil.Process(pid)
    if not process.is_running():
       return None, "process {} is not running".format(pid)

    return psutil.Process(pid), ""

####
def execute(cmd):
    cmdJSON, hv = cmd_md5(cmd)
    # os.makedirs("log", exist_ok=True)
    # logFile = os.path.join("log", hv + ".log")

    logFile, pidFile = hv + ".log", hv + ".pid"
    # if os.path.exists(logFile):
    if os.path.exists(pidFile):
        _, err = process_from_pidFile(pidFile)
        if err == "":
            detail = "process is running"
            reportTask("run", "conflict", detail)
            return detail
        # os.remove(pidFile)

    ####
    # notExits = not os.path.exists(logFile)
    # log_file = open(logFile, 'a+')
    # if notExits: log_file.write(cmdJSON + "\n")
    log_file = open(logFile, 'w')
    t1 = datetime.now()
    log_file.write(">>> {} start: {}\n\n".format(rfc3339ms(t1), cmdJSON))
    log_file.flush()

    proc = subprocess.Popen(cmd, stdout=log_file, stderr=log_file)

    time.sleep(1.5)
    if not proc.poll() is None:
        log_file.close()
        reportTask("run", "failed")
        return -100

    with open(pidFile, 'w') as pid_file:
        pid_file.write("{}\n\n{}\n{}\n".format(proc.pid, rfc3339ms(t1), cmdJSON))

    log_file.flush()

    ####
    status = proc.wait()
    t2 = datetime.now()
    log_file.write("\n### {} {}, exit status: {}\n".format(rfc3339ms(t2), t2-t1, status))
    os.remove(pidFile)
    log_file.close()

    if status == 0:
        reportTask("run", "ok", status)
    else:
        reportTask("run", "failed", status)

    return status

def status(cmd):
    _, hv = cmd_md5(cmd)
    pidFile = hv + ".pid"
    if not os.path.exists(pidFile): return "exit", ""

    process, err = get_process(cmd)
    if err != "": return None, err

    return  process.status(), ""

def info(cmd):
    process, err = get_process(cmd)
    if err != "": return None, err

    return json.dumps(process.as_dict()), ""

def kill(cmd):
    process, err = get_process(cmd)
    if err != "": return err

    for p in process.children(): p.kill()
    if process.is_running(): process.kill()

    return 0

def suspend(cmd):
    process, err = get_process(cmd)
    if err != "": return err

    if process.status() == "stopped": return 0

    process.suspend()
    for p in process.children(): p.suspend()

    return 0

def resume(cmd):
    process, err = get_process(cmd)
    if err != "": return err

    # if process.status() == "stopped":
    for p in process.children(): p.resume()
    process.resume()

    return 0

if __name__ == '__main__':
    if len(sys.argv) < 2:
        sys.exit("""execute.py <subcommand> <args...>
subcommand: run, kill, suspend, resume, md5, status, info
""")

    call, cmd = sys.argv[1], sys.argv[2:]

    if call == "run":
        sys.exit(execute(cmd))
    elif call == "kill":
        sys.exit(kill(cmd))
    elif call == "suspend":
        sys.exit(suspend(cmd))
    elif call == "resume":
        sys.exit(resume(cmd))
    elif call == "md5":
        _, cmdMD5 = cmd_md5(cmd)
        print(cmdMD5)
    elif call == "status":
        result, err = status(cmd)
        if err != "": sys.exit(err)
        print(result)
    elif call == "info":
        result, err = info(cmd)
        if err != "": sys.exit(err)
        print(result)
    else:
        sys.exit("unkonwn subcommand: {}".format(call))
