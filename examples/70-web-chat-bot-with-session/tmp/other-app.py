import streamlit as st
import requests
import os
from datetime import datetime

# Configuration of the Streamlit page
st.set_page_config(page_title="TechFusion 2025", page_icon="🚀")

# Hide the Deploy button
st.markdown("""
    <style>
    .stDeployButton {
        visibility: hidden;
    }
    </style>
    """, unsafe_allow_html=True)

# Initialisation of the session state
if "messages" not in st.session_state:
    st.session_state.messages = []

# Handle the reset of the input key
if "input_key" not in st.session_state:
    st.session_state.input_key = 0

# Initialize session ID in session state
if "session_id" not in st.session_state:
    st.session_state.session_id = "default"

# Backend URL (the nodejs server)
BACKEND_SERVICE_URL = os.environ.get('BACKEND_SERVICE_URL', 'http://backend:5050')

def stream_response(message, session_id):
    """Stream the message response from the backend"""
    try:
        with requests.post(
            BACKEND_SERVICE_URL+"/chat",
            json={"message": message, "sessionId": session_id},
            headers={"Content-Type": "application/json"},
            stream=True
        ) as response:
            # Create a placeholder for the streaming response
            response_placeholder = st.empty()
            full_response = ""
            
            # Stream the response chunks
            for chunk in response.iter_content(chunk_size=1024, decode_unicode=True):
                if chunk:
                    chunk_text = chunk.decode('utf-8') if isinstance(chunk, bytes) else chunk
                    full_response += chunk_text
                    # Update the placeholder with the accumulated response
                    response_placeholder.markdown(full_response)
            
            return full_response
    except requests.exceptions.RequestException as e:
        error_msg = f"😡 Connection error: {str(e)}"
        st.error(error_msg)
        return error_msg

def clear_conversation_history(session_id):
    """Clear the conversation history on the server"""
    try:
        response = requests.post(
            f"{BACKEND_SERVICE_URL}/clear-history",
            json={"sessionId": session_id},
            headers={"Content-Type": "application/json"}
        )
        if response.status_code == 200:
            st.session_state.messages = []  # Clear local messages too
            st.success("✨ Conversation history cleared!")
        else:
            st.error("Failed to clear conversation history")
    except requests.exceptions.RequestException as e:
        st.error(f"Error clearing history: {str(e)}")

def increment_input_key():
    """Increment the input key to reset the input field"""
    st.session_state.input_key += 1

# Page title
st.title("🚀 TechFusion 2025")
st.header("Where Web3 Meets AI 🤖")

# Session ID input
session_id = st.text_input(
    "🔑 Session ID:",
    value=st.session_state.session_id,
    help="Enter a unique session ID to maintain conversation context"
)
st.session_state.session_id = session_id

# Form to send a message
with st.form(key=f"message_form_{st.session_state.input_key}"):
    message = st.text_input("📝 Your message:", key=f"input_{st.session_state.input_key}")
    col1, col2, col3 = st.columns([1, 2.5, 1])
    with col1:
        submit_button = st.form_submit_button("Send ✨")
    with col2:
        # Empty column for spacing
        pass
    with col3:
        clear_button = st.form_submit_button("Clear History 🗑️")

# Handle the clear history button
if clear_button:
    clear_conversation_history(st.session_state.session_id)
    st.rerun()

# Handle the message submission
if submit_button and message and len(message.strip()) > 0:
    # Add the message to the history
    st.session_state.messages.append({
        "role": "user",
        "content": message,
        "time": datetime.now(),
        "session_id": st.session_state.session_id
    })
    
    # Stream the response from the backend
    response = stream_response(message, st.session_state.session_id)
    
    # Add the response to the history
    st.session_state.messages.append({
        "role": "assistant",
        "content": response,
        "time": datetime.now(),
        "session_id": st.session_state.session_id
    })
    
    # Reset the input field
    increment_input_key()
    st.rerun()

# Display the messages history
st.write("### Messages history")
for msg in reversed(st.session_state.messages):
    with st.container():
        if msg["role"] == "user":
            st.info(f"🤓 You ({msg['time'].strftime('%H:%M')}) - Session: {msg['session_id']}")
            st.write(msg["content"])
        else:
            st.success(f"🤖 Assistant ({msg['time'].strftime('%H:%M')}) - Session: {msg['session_id']}")
            st.write(msg["content"])