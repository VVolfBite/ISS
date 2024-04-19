import os

def find_performance_vals(root_dir, performance_metric):
    for dirpath, dirnames, filenames in os.walk(root_dir):
        for dirname in dirnames:
            if dirname.startswith('remote-') or dirname.startswith('local-'):
                experiment_output_dir = os.path.join(dirpath, dirname, 'experiment-output')
                if os.path.exists(experiment_output_dir):
                    print(f"At {dirname}:")
                    for exp_dir in os.listdir(experiment_output_dir):
                        exp_output_path = os.path.join(experiment_output_dir, exp_dir)
                        performance_val_path = os.path.join(exp_output_path, performance_metric)
                        if os.path.exists(performance_val_path):
                            with open(performance_val_path, 'r') as f:
                                performance_val = f.read().strip()
                            print(f"{exp_dir}/{performance_metric}: {performance_val}")

root_directory = '/opt/gopath/src/github.com/hyperledger-labs/mirbft/deployment/deployment-data/'
performance_metric = 'throughput-raw.val'  # Specify the performance metric here
find_performance_vals(root_directory, performance_metric)
