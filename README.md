# AI Game Backend

A Go-based game backend featuring AI-powered character interactions and natural language conversation systems.

## Overview

This project was developed as the backend for an AI-enhanced game where players could have natural conversations with NPCs through text. While the game didn't reach commercial release, the technical implementation provided extensive experience in AI integration, natural language processing, and Unity-backend communication.

## Features

- **Natural Language Conversations**: Players can chat with AI characters using normal text
- **SMS Integration**: Text game characters from your phone (planned feature)
- **Intelligent NPCs**: Characters respond contextually and maintain conversation history
- **WebSocket Communication**: Real-time text-based interaction between Unity client and backend
- **Character Personality Systems**: Each NPC has distinct conversational patterns and knowledge
- **Conversation Persistence**: Chat history and character relationships stored long-term

## Tech Stack

- **Go**: Core backend language
- **WebSockets**: Real-time text communication
- **PostgreSQL**: Character data and conversation history
- **Unity**: Frontend game client
- **AI Integration**: Natural language processing and character intelligence

## Architecture

The backend enables natural conversations with game characters by:
1. Receiving text input from Unity client via WebSocket
2. Processing player messages through AI language models
3. Generating contextually appropriate character responses
4. Maintaining conversation history and character personality consistency
5. Streaming responses back to Unity for real-time chat experience

## Key Technical Challenges Solved

- **Natural Language Processing**: Implementing conversational AI that feels natural and engaging
- **Character Consistency**: Maintaining distinct personalities and knowledge for each NPC
- **Conversation Context**: Preserving chat history and relationship development over time
- **Real-time Communication**: Seamless text streaming between Unity and Go backend
- **Performance Optimization**: Handling AI processing while maintaining responsive chat experience

## What I Learned

- **Conversational AI**: Implementing natural language processing in interactive systems
- **WebSocket Communication**: Real-time text streaming between Unity and Go
- **Character AI Design**: Creating distinct personalities and consistent conversational patterns
- **Unity-Backend Integration**: Seamless communication between game client and server
- **Context Management**: Preserving conversation history and relationship dynamics
- **Go Concurrency**: Handling multiple simultaneous conversations efficiently

## Business Context

Originally designed as a prototype for an AI-driven life-sim game where players could have meaningful conversations with any character. The backend successfully demonstrated natural language interaction with AI NPCs, but market research revealed challenges in content creation scalability and player retention.

The codebase represents the prototype for the conversational game backend that showcases advanced AI integration, real-time communication, and Unity game development collaboration.
