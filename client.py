import json
import sys

def main():
    # Read JSON from STDIN
    input_json = sys.stdin.read()

    # Parse the JSON into a dictionary
    try:
        data = json.loads(input_json)
    except json.JSONDecodeError as e:
        print(f"Error decoding JSON: {e}")
        return

    # Extract the 'report' list from the data
    report_list = data.get("report", [])

    # Sort the byte map by keys (converted to integers) and concatenate values
    sorted_bytes = sorted(report_list, key=lambda kv: int(kv["key"]))
    hex_string = ''.join(item["value"] for item in sorted_bytes)

    # Print the reconstructed hex string
    print(hex_string)

if __name__ == "__main__":
    main()
