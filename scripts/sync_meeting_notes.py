import os
import io
from google.oauth2 import service_account
from googleapiclient.discovery import build
from googleapiclient.http import MediaIoBaseDownload
from git import Repo
from git.exc import GitCommandError

# --- Configuration ---
# Environment variables from GitHub Actions
DOC_ID = os.environ.get('MEETING_NOTES_SYNC_DOC_ID')
TARGET_FILE_PATH = os.environ.get('MEETING_NOTES_SYNC_TARGET_FILE_PATH')

# --- Google Docs API Authentication and Download ---
def download_markdown_from_drive(doc_id):
    SCOPES = ['https://www.googleapis.com/auth/drive.readonly']
    SERVICE_ACCOUNT_FILE = 'service_account.json'

    creds = service_account.Credentials.from_service_account_file(
        SERVICE_ACCOUNT_FILE, scopes=SCOPES)
    
    # Build Drive API client
    service = build('drive', 'v3', credentials=creds)

    # Export MIME type for Docs/Drive API: 'text/markdown'
    request = service.files().export_media(
        fileId=doc_id, 
        mimeType='text/markdown'
    )
    
    # Use BytesIO to receive document content
    fh = io.BytesIO()
    downloader = MediaIoBaseDownload(fh, request)
    done = False
    while done is False:
        status, done = downloader.next_chunk()
        # Print download progress (optional)
        # print(f"Download {int(status.progress() * 100)}%.")

    # Return byte content, ensuring UTF-8 encoding
    # Compared to Apps Script strings, Python's bytes.decode('utf-8') is more reliable
    return fh.getvalue().decode('utf-8')

# --- Git Operations: Write and Stage File ---
def write_and_stage_file(content):
    # Ensure target directory exists
    os.makedirs(os.path.dirname(TARGET_FILE_PATH), exist_ok=True)
    
    # Write document, using utf-8 encoding
    with open(TARGET_FILE_PATH, 'w', encoding='utf-8') as f:
        f.write(content)

    # Initialize Git repository object
    repo = Repo('.')
    # Stage changes
    repo.index.add([TARGET_FILE_PATH])
    print(f"File {TARGET_FILE_PATH} written and staged successfully.")

# --- Main Execution Flow ---
if __name__ == '__main__':
    try:
        # 1. Download Markdown content from Google Drive
        markdown_content = download_markdown_from_drive(DOC_ID)
        
        # 2. Write content to local working directory and stage
        write_and_stage_file(markdown_content)

    except Exception as e:
        print(f"An error occurred during sync: {e}")
        exit(1)

    finally:
        # Delete key document at the end of the task to prevent inclusion in subsequent steps (e.g., git status)
        if os.path.exists("service_account.json"):
            os.remove("service_account.json")
            print("Cleaned up service_account.json.")
