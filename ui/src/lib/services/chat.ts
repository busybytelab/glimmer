import { getAuthToken, getCurrentUserId } from '$lib/auth';

export type Chat = {
    id: string;
    title: string;
    systemPrompt: string;
    model?: string;
    createdAt: Date;
    updatedAt: Date;
    userID: string;
    archived: boolean;
    lastMessage?: ChatMessage | null; // Optional last message for preview
};

export type ChatMessage = {
    id: string;
    chatID: string;
    role: 'user' | 'assistant' | 'system';
    content: string;
    createdAt: Date;
};

export type ChatWithMessages = Chat & {
    messages: ChatMessage[];
};

class ChatService {
    private static instance: ChatService;
    private baseUrl = '/api/collections';

    private constructor() {}

    public static getInstance(): ChatService {
        if (!ChatService.instance) {
            ChatService.instance = new ChatService();
        }
        return ChatService.instance;
    }

    /**
     * Fetches all chats for the current user
     */
    public async getChats(): Promise<Chat[]> {
        const authToken = getAuthToken();
        if (!authToken) {
            throw new Error('Please log in again.');
        }

        const response = await fetch(`${this.baseUrl}/chats/records?sort=-created`, {
            headers: {
                'Authorization': `Bearer ${authToken}`
            }
        });

        if (!response.ok) {
            throw new Error(`Error: ${response.status} ${response.statusText}`);
        }

        const data = await response.json();
        
        // Transform PocketBase response to our Chat type
        return data.items.map((item: any) => this.transformChatResponse(item));
    }

    /**
     * Fetches a single chat with all its messages
     */
    public async getChat(chatId: string): Promise<ChatWithMessages> {
        const authToken = getAuthToken();
        if (!authToken) {
            throw new Error('Please log in again.');
        }

        // Fetch the chat details
        const chatResponse = await fetch(`${this.baseUrl}/chats/records/${chatId}`, {
            headers: {
                'Authorization': `Bearer ${authToken}`
            }
        });

        if (!chatResponse.ok) {
            throw new Error(`Error: ${chatResponse.status} ${chatResponse.statusText}`);
        }

        const chatData = await chatResponse.json();
        const chat = this.transformChatResponse(chatData);

        // Fetch the chat messages
        const messagesResponse = await fetch(
            `${this.baseUrl}/chat_items/records?filter=(chat="${chatId}")&sort=created`, {
            headers: {
                'Authorization': `Bearer ${authToken}`
            }
        });

        if (!messagesResponse.ok) {
            const errorData = await messagesResponse.json().catch(() => ({}));
            console.error('Error fetching chat messages:', errorData);
            throw new Error(`Error: ${messagesResponse.status} ${messagesResponse.statusText}`);
        }

        const messagesData = await messagesResponse.json();
        const messages = messagesData.items.map((item: any) => this.transformMessageResponse(item));

        // Return combined data
        return {
            ...chat,
            messages
        };
    }

    /**
     * Creates a new chat
     */
    public async createChat(systemPrompt: string, model?: string): Promise<Chat> {
        const authToken = getAuthToken();
        if (!authToken) {
            throw new Error('Please log in again.');
        }

        const userId = getCurrentUserId();
        if (!userId) {
            throw new Error('User ID not found. Please log in again.');
        }

        const requestBody: Record<string, any> = {
            system_prompt: systemPrompt,
            label: 'New Chat', // Default label (title) - will be updated after first message
            user: userId // Add the user ID to the request
        };

        if (model) {
            requestBody.model = model;
        }

        const response = await fetch(`${this.baseUrl}/chats/records`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${authToken}`
            },
            body: JSON.stringify(requestBody)
        });

        if (!response.ok) {
            const errorData = await response.json().catch(() => ({}));
            console.error('Error creating chat:', errorData);
            throw new Error(`Error: ${response.status} ${response.statusText}`);
        }

        const data = await response.json();
        return this.transformChatResponse(data);
    }

    /**
     * Updates a chat's title
     */
    public async updateChatTitle(chatId: string, title: string): Promise<Chat> {
        const authToken = getAuthToken();
        if (!authToken) {
            throw new Error('Please log in again.');
        }

        const response = await fetch(`${this.baseUrl}/chats/records/${chatId}`, {
            method: 'PATCH',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${authToken}`
            },
            body: JSON.stringify({
                label: title // Use label field name instead of title
            })
        });

        if (!response.ok) {
            throw new Error(`Error: ${response.status} ${response.statusText}`);
        }

        const data = await response.json();
        return this.transformChatResponse(data);
    }

    /**
     * Deletes a chat
     */
    public async deleteChat(chatId: string): Promise<void> {
        const authToken = getAuthToken();
        if (!authToken) {
            throw new Error('Please log in again.');
        }

        const response = await fetch(`${this.baseUrl}/chats/records/${chatId}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${authToken}`
            }
        });

        if (!response.ok) {
            throw new Error(`Error: ${response.status} ${response.statusText}`);
        }
    }

    /**
     * Fetches all chats for the current user with their last messages
     * @param includeArchived Whether to include archived chats, defaults to false
     */
    public async getChatsWithLastMessages(includeArchived: boolean = false): Promise<Chat[]> {
        const authToken = getAuthToken();
        if (!authToken) {
            throw new Error('Please log in again.');
        }

        const userId = getCurrentUserId();
        if (!userId) {
            throw new Error('User ID not found. Please log in again.');
        }

        // Build filter to exclude archived chats by default
        let filter = `user="${userId}"`;
        if (!includeArchived) {
            filter += ` && (archived=false || archived=null)`;
        }

        // Fetch chats with filter
        const response = await fetch(
            `${this.baseUrl}/chats/records?filter=${encodeURIComponent(filter)}&sort=-updated`, {
            headers: {
                'Authorization': `Bearer ${authToken}`
            }
        });

        if (!response.ok) {
            throw new Error(`Error: ${response.status} ${response.statusText}`);
        }

        const data = await response.json();
        const chats = data.items.map((item: any) => this.transformChatResponse(item));
        
        // If no chats, return empty array
        if (chats.length === 0) {
            return [];
        }
        
        // Fetch last messages for each chat
        const chatIds = chats.map((chat: Chat) => chat.id);
        const lastMessagesPromises = chatIds.map((chatId: string) => this.getLastMessage(chatId));
        
        try {
            const lastMessages = await Promise.all(lastMessagesPromises);
            
            // Associate last messages with their chats
            return chats.map((chat: Chat, index: number) => ({
                ...chat,
                lastMessage: lastMessages[index]
            }));
        } catch (error) {
            console.error('Error fetching last messages:', error);
            // Return chats without last messages on error
            return chats;
        }
    }
    
    /**
     * Fetches the last message for a chat
     */
    public async getLastMessage(chatId: string): Promise<ChatMessage | null> {
        const authToken = getAuthToken();
        if (!authToken) {
            throw new Error('Please log in again.');
        }

        try {
            const response = await fetch(
                `${this.baseUrl}/chat_items/records?filter=(chat="${chatId}")&sort=-created&limit=1`, {
                headers: {
                    'Authorization': `Bearer ${authToken}`
                }
            });

            if (!response.ok) {
                throw new Error(`Error: ${response.status} ${response.statusText}`);
            }

            const data = await response.json();
            
            if (data.items.length === 0) {
                return null;
            }
            
            return this.transformMessageResponse(data.items[0]);
        } catch (error) {
            console.error(`Error fetching last message for chat ${chatId}:`, error);
            return null;
        }
    }

    /**
     * Helper method to transform PocketBase response to our Chat type
     */
    private transformChatResponse(data: any): Chat {
        // Use label field from the API response (not title)
        // Ensure we have a valid title (never undefined)
        const title = data.label || 'New Chat';
        
        return {
            id: data.id,
            title: title,
            systemPrompt: data.system_prompt,
            model: data.model,
            createdAt: new Date(data.created),
            updatedAt: new Date(data.updated),
            userID: data.user,
            archived: data.archived || false
        };
    }

    /**
     * Helper method to transform PocketBase response to our ChatMessage type
     */
    private transformMessageResponse(data: any): ChatMessage {
        return {
            id: data.id,
            chatID: data.chat,
            role: data.role,
            content: data.content,
            createdAt: new Date(data.created)
        };
    }

    /**
     * Generates a title for a chat based on the user's first message
     * @param message The user message to base the title on
     * @returns A title that summarizes the message (truncated if needed)
     */
    private generateChatTitle(message: string): string {
        // Clean up the message
        const cleanMessage = message.trim();
        
        // Return default title if message is empty
        if (!cleanMessage) return 'New Chat';
        
        // Use first 5 words or 50 characters, whichever is shorter
        const words = cleanMessage.split(/\s+/);
        if (words.length <= 5) {
            return cleanMessage.length <= 50 ? cleanMessage : cleanMessage.substring(0, 47) + '...';
        }
        
        const title = words.slice(0, 5).join(' ');
        return title.length <= 50 ? title : title.substring(0, 47) + '...';
    }

    /**
     * Adds a message to a chat and updates the chat title if it's the first user message
     */
    public async addMessageToChat(chatId: string, message: string, role: 'user' | 'assistant' | 'system'): Promise<ChatMessage> {
        const authToken = getAuthToken();
        if (!authToken) {
            throw new Error('Please log in again.');
        }
        
        // Create message
        const messageResponse = await fetch(`${this.baseUrl}/chat_items/records`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${authToken}`
            },
            body: JSON.stringify({
                chat: chatId,
                role: role,
                content: message
            })
        });
        
        if (!messageResponse.ok) {
            throw new Error(`Error: ${messageResponse.status} ${messageResponse.statusText}`);
        }
        
        const messageData = await messageResponse.json();
        
        // If this is a user message, check if we need to update the chat title
        if (role === 'user') {
            try {
                // Get all messages for this chat
                const messagesResponse = await fetch(
                    `${this.baseUrl}/chat_items/records?filter=(chat="${chatId}")&sort=created`, {
                    headers: {
                        'Authorization': `Bearer ${authToken}`
                    }
                });
                
                if (messagesResponse.ok) {
                    const messagesData = await messagesResponse.json();
                    
                    // If this is the first user message, update the chat title
                    if (messagesData.items.filter((item: any) => item.role === 'user').length === 1) {
                        const title = this.generateChatTitle(message);
                        await this.updateChatTitle(chatId, title);
                    }
                }
            } catch (error) {
                console.error('Error updating chat title:', error);
                // Non-critical error, continue without updating title
            }
        }
        
        return this.transformMessageResponse(messageData);
    }

    /**
     * Searches chats by title and message content
     * @param query The search query string
     * @returns Array of chats matching the query
     */
    public async searchChats(query: string): Promise<Chat[]> {
        const authToken = getAuthToken();
        if (!authToken) {
            throw new Error('Please log in again.');
        }
        
        // If query is empty, return all chats
        if (!query.trim()) {
            return this.getChatsWithLastMessages();
        }
        
        try {
            const userId = getCurrentUserId();
            if (!userId) {
                throw new Error('User ID not found. Please log in again.');
            }
            
            const trimmedQuery = query.trim();
            
            // Run both searches in parallel for efficiency
            const [labelResults, contentResults] = await Promise.all([
                // 1. Search chats by label
                this.searchChatsByLabel(trimmedQuery, userId, authToken),
                
                // 2. Search chat messages by content (if query is long enough)
                trimmedQuery.length >= 2 ? this.searchChatsByContent(trimmedQuery, userId, authToken) : []
            ]);
            
            // Combine results, removing duplicates by ID
            const chatMap = new Map<string, Chat>();
            
            // Add label search results first
            labelResults.forEach(chat => {
                chatMap.set(chat.id, chat);
            });
            
            // Add content search results, keeping existing ones if already found
            contentResults.forEach(chat => {
                if (!chatMap.has(chat.id)) {
                    chatMap.set(chat.id, chat);
                }
            });
            
            // Convert back to array and sort by updated date
            const combinedResults = Array.from(chatMap.values())
                .sort((a, b) => b.updatedAt.getTime() - a.updatedAt.getTime());
            
            return combinedResults;
        } catch (error) {
            console.error('Error searching chats:', error);
            throw error;
        }
    }
    
    /**
     * Helper method to search chats by label
     */
    private async searchChatsByLabel(query: string, userId: string, authToken: string): Promise<Chat[]> {        
        // Include archived chats in search results
        const labelFilter = `label~"${query}" && user="${userId}"`;
        
        try {
            const response = await fetch(
                `${this.baseUrl}/chats/records?filter=${encodeURIComponent(labelFilter)}&sort=-updated`, {
                headers: {
                    'Authorization': `Bearer ${authToken}`
                }
            });
            
            if (!response.ok) {
                console.error("Label search failed:", response.status, response.statusText);
                return [];
            }
            
            const data = await response.json();
            const chats = data.items.map((item: any) => this.transformChatResponse(item));
            
            // Fetch last messages for each chat
            const chatIds = chats.map((chat: Chat) => chat.id);
            const lastMessagesPromises = chatIds.map((chatId: string) => this.getLastMessage(chatId));
            const lastMessages = await Promise.all(lastMessagesPromises);
            
            // Associate last messages with their chats
            return chats.map((chat: Chat, index: number) => ({
                ...chat,
                lastMessage: lastMessages[index]
            }));
        } catch (error) {
            console.error("Error in label search:", error);
            return [];
        }
    }
    
    /**
     * Helper method to search chats by message content
     */
    private async searchChatsByContent(query: string, userId: string, authToken: string): Promise<Chat[]> {        
        try {
            // Search for messages containing the query
            const contentFilter = `content~"${query}"`;
            
            const contentResponse = await fetch(
                `${this.baseUrl}/chat_items/records?filter=${encodeURIComponent(contentFilter)}&sort=-created`, {
                headers: {
                    'Authorization': `Bearer ${authToken}`
                }
            });
            
            if (!contentResponse.ok) {
                console.error("Content search failed:", contentResponse.status, contentResponse.statusText);
                return [];
            }
            
            const contentData = await contentResponse.json();
            
            // Get unique chat IDs from the messages
            const chatIdsSet = new Set<string>();
            contentData.items.forEach((item: any) => {
                chatIdsSet.add(item.chat);
            });
            
            if (chatIdsSet.size === 0) {
                return [];
            }
            
            // Build a filter to get all matching chats at once
            const chatIdsArray = Array.from(chatIdsSet);
            const chatsFilter = chatIdsArray.map(id => `id="${id}"`).join(" || ");
            const userFilter = `user="${userId}"`;
            const combinedFilter = `(${chatsFilter}) && ${userFilter}`;
            
            // Fetch the chats
            const chatsResponse = await fetch(
                `${this.baseUrl}/chats/records?filter=${encodeURIComponent(combinedFilter)}&sort=-updated`, {
                headers: {
                    'Authorization': `Bearer ${authToken}`
                }
            });
            
            if (!chatsResponse.ok) {
                console.error("Chats fetch failed:", chatsResponse.status, chatsResponse.statusText);
                return [];
            }
            
            const chatsData = await chatsResponse.json();
            const chats = chatsData.items.map((item: any) => this.transformChatResponse(item));
            
            // Create a map of messages by chat ID for quick lookup
            const messagesByChatId = new Map<string, any[]>();
            contentData.items.forEach((message: any) => {
                const chatId = message.chat;
                if (!messagesByChatId.has(chatId)) {
                    messagesByChatId.set(chatId, []);
                }
                messagesByChatId.get(chatId)?.push(message);
            });
            
            // Associate the most recent matching message with each chat
            return chats.map((chat: Chat) => {
                const chatMessages = messagesByChatId.get(chat.id) || [];
                // Sort messages by date, newest first
                chatMessages.sort((a: any, b: any) => new Date(b.created).getTime() - new Date(a.created).getTime());
                
                return {
                    ...chat,
                    lastMessage: chatMessages.length > 0 ? this.transformMessageResponse(chatMessages[0]) : null
                };
            });
        } catch (error) {
            console.error("Error in content search:", error);
            return [];
        }
    }

    /**
     * Updates a chat's archived status
     */
    public async updateChatArchiveStatus(chatId: string, archived: boolean): Promise<Chat> {
        const authToken = getAuthToken();
        if (!authToken) {
            throw new Error('Please log in again.');
        }

        const response = await fetch(`${this.baseUrl}/chats/records/${chatId}`, {
            method: 'PATCH',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${authToken}`
            },
            body: JSON.stringify({
                archived: archived
            })
        });

        if (!response.ok) {
            throw new Error(`Error: ${response.status} ${response.statusText}`);
        }

        const data = await response.json();
        return this.transformChatResponse(data);
    }
}

// Export a singleton instance
export const chatService = ChatService.getInstance(); 