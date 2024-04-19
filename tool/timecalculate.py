import re

def parse_log(log_file):
    total_execution_time = 0
    total_calls = 0

    with open(log_file, 'r') as file:
        for line in file:
            match = re.search(r'Execution time: (\d+).', line)
            if match:
                execution_time = int(match.group(1))
                total_execution_time += execution_time
                total_calls += 1

    if total_calls > 0:
        average_execution_time = total_execution_time / total_calls
        print(f"Total calls: {total_calls}")
        print(f"Total execution time: {total_execution_time} microseconds")
        print(f"Average execution time: {average_execution_time} microseconds")
    else:
        print("No matching lines found in the log file")

if __name__ == "__main__":
    log_file_2 = "/opt/gopath/src/github.com/hyperledger-labs/mirbft/deployment/deployment-data/local-0000/experiment-output/0001/slave-001/peer.log"  # 修改成你的日志文件路径
    parse_log(log_file_2)
