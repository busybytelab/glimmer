import { writable } from 'svelte/store';
import type { Chat, ChatWithMessages } from '$lib/services/chat';
import { chatService } from '$lib/services/chat';

// Store for chat list
export const chatListStore = writable<{
    chats: Chat[];
    loading: boolean;
    error: string | null;
}>({
    chats: [],
    loading: false,
    error: null
});

// Store for active chat
export const activeChatStore = writable<{
    chat: ChatWithMessages | null;
    loading: boolean;
    error: string | null;
}>({
    chat: null,
    loading: false,
    error: null
});

// Actions to load chats
export async function loadChats(): Promise<void> {
    chatListStore.update(state => ({ ...state, loading: true, error: null }));
    
    try {
        // Use the enhanced method to get chats with their last messages
        const chats = await chatService.getChatsWithLastMessages();
        chatListStore.update(state => ({ 
            ...state, 
            chats, 
            loading: false 
        }));
    } catch (error) {
        console.error('Error loading chats:', error);
        chatListStore.update(state => ({ 
            ...state, 
            loading: false, 
            error: error instanceof Error ? error.message : 'Failed to load chats' 
        }));
    }
}

// Load a specific chat
export async function loadChat(chatId: string): Promise<void> {
    activeChatStore.update(state => ({ ...state, loading: true, error: null }));
    
    try {
        const chat = await chatService.getChat(chatId);
        activeChatStore.update(state => ({ 
            ...state, 
            chat, 
            loading: false 
        }));
    } catch (error) {
        console.error('Error loading chat:', error);
        activeChatStore.update(state => ({ 
            ...state, 
            loading: false, 
            error: error instanceof Error ? error.message : 'Failed to load chat'
        }));
    }
}

// Create a new chat
export async function createNewChat(systemPrompt: string, model?: string): Promise<string> {
    try {
        const newChat = await chatService.createChat(systemPrompt, model);
        
        // Update the chat list store with the new chat
        chatListStore.update(state => ({
            ...state,
            chats: [newChat, ...state.chats]
        }));
        
        return newChat.id;
    } catch (error) {
        console.error('Error creating chat:', error);
        throw error;
    }
}

// Update chat title
export async function updateChatTitle(chatId: string, title: string): Promise<void> {
    try {
        const updatedChat = await chatService.updateChatTitle(chatId, title);
        
        // Update in chat list
        chatListStore.update(state => ({
            ...state,
            chats: state.chats.map(chat => 
                chat.id === chatId ? updatedChat : chat
            )
        }));
        
        // Update in active chat if needed
        activeChatStore.update(state => {
            if (state.chat && state.chat.id === chatId) {
                return {
                    ...state,
                    chat: {
                        ...state.chat,
                        title
                    }
                };
            }
            return state;
        });
    } catch (error) {
        console.error('Error updating chat title:', error);
        throw error;
    }
}

// Delete a chat
export async function deleteChat(chatId: string): Promise<void> {
    try {
        await chatService.deleteChat(chatId);
        
        // Remove from chat list
        chatListStore.update(state => ({
            ...state,
            chats: state.chats.filter(chat => chat.id !== chatId)
        }));
        
        // Clear active chat if it was the one deleted
        activeChatStore.update(state => {
            if (state.chat && state.chat.id === chatId) {
                return {
                    ...state,
                    chat: null
                };
            }
            return state;
        });
    } catch (error) {
        console.error('Error deleting chat:', error);
        throw error;
    }
}

// Search chats based on query
export async function searchChats(query: string): Promise<void> {
    // If query is empty, just load all chats
    if (!query.trim()) {
        return loadChats();
    }

    chatListStore.update(state => ({ ...state, loading: true, error: null }));
    
    try {
        const chats = await chatService.searchChats(query);
        chatListStore.update(state => ({ 
            ...state, 
            chats, 
            loading: false 
        }));
    } catch (error) {
        console.error('Error searching chats:', error);
        chatListStore.update(state => ({ 
            ...state, 
            loading: false, 
            error: error instanceof Error ? error.message : 'Failed to search chats' 
        }));
    }
}

// Archive/unarchive a chat
export async function archiveChat(chatId: string, archived: boolean = true): Promise<void> {
    try {
        const updatedChat = await chatService.updateChatArchiveStatus(chatId, archived);
        
        // Update in chat list - if archiving, filter it out, if unarchiving, update it
        chatListStore.update(state => {
            // If chat is being archived and we're in normal mode (not showing archived)
            // then remove it from the list
            if (archived) {
                return {
                    ...state,
                    chats: state.chats.filter(chat => chat.id !== chatId)
                };
            }
            
            // If chat is being unarchived, update its status
            return {
                ...state,
                chats: state.chats.map(chat => 
                    chat.id === chatId ? updatedChat : chat
                )
            };
        });
        
        // Update in active chat if needed
        activeChatStore.update(state => {
            if (state.chat && state.chat.id === chatId) {
                return {
                    ...state,
                    chat: {
                        ...state.chat,
                        archived
                    }
                };
            }
            return state;
        });
    } catch (error) {
        console.error('Error updating chat archive status:', error);
        throw error;
    }
} 