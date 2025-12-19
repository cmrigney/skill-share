#!/usr/bin/env python3
"""
Example skill script demonstrating code execution.

Skills can include executable scripts that Claude runs via bash.
Only the script output (not the code) enters Claude's context.
"""

def main():
    print("Example Skill Script")
    print("=" * 40)
    print()
    print("This script demonstrates how skills can")
    print("include executable code.")
    print()
    print("Benefits:")
    print("- Deterministic operations")
    print("- Efficient (code doesn't consume tokens)")
    print("- Reliable and testable")
    print()
    print("Script executed successfully!")

if __name__ == "__main__":
    main()
