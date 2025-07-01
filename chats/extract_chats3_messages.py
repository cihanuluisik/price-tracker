#!/usr/bin/env python3
"""
Extract user messages from price-tracker-chats3.md file only.
This script extracts only the user messages without any additional formatting.
"""

import re
from pathlib import Path

def extract_user_messages_from_chats3():
    """
    Extract user messages from price-tracker-chats3.md file.
    
    Returns:
        List of user messages
    """
    file_path = Path("chats/price-tracker-chats3.md")
    
    if not file_path.exists():
        print(f"File {file_path} does not exist")
        return []
    
    user_messages = []
    
    try:
        with open(file_path, 'r', encoding='utf-8') as file:
            content = file.read()
            
        # Split content by the separator pattern
        sections = re.split(r'---\s*\n', content)
        
        for section in sections:
            section = section.strip()
            if not section:
                continue
                
            # Check if this section starts with "**User**"
            if section.startswith('**User**'):
                # Extract the message content after "**User**"
                message_start = section.find('**User**') + len('**User**')
                message_content = section[message_start:].strip()
                
                # Clean up the message (remove leading/trailing whitespace and newlines)
                message_content = re.sub(r'^\s*\n+', '', message_content)
                message_content = re.sub(r'\n+\s*$', '', message_content)
                
                if message_content:
                    user_messages.append(message_content)
                    
    except Exception as e:
        print(f"Error processing {file_path}: {e}")
        
    return user_messages

def save_messages_only(messages, output_file="chats3_user_messages.txt"):
    """
    Save messages to a simple text file, one message per line.
    
    Args:
        messages: List of user messages
        output_file: Output file path
    """
    with open(output_file, 'w', encoding='utf-8') as f:
        for message in messages:
            f.write(f"{message}\n\n")
    
    print(f"Messages saved to {output_file}")
    print(f"Total messages: {len(messages)}")

def main():
    """Main function to extract messages from chats3."""
    print("Extracting user messages from price-tracker-chats3.md...")
    
    # Extract messages
    messages = extract_user_messages_from_chats3()
    
    if messages:
        # Save messages to file
        save_messages_only(messages)
        
        # Also print to console
        print("\n" + "="*50)
        print("EXTRACTED MESSAGES")
        print("="*50)
        
        for i, message in enumerate(messages, 1):
            print(f"\n--- Message {i} ---")
            print(message)
            print("-" * 50)
    else:
        print("No messages found.")

if __name__ == "__main__":
    main() 