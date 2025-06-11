import pb from '$lib/pocketbase';

class AnswersService {
    /**
     * Evaluates a user's answer for a practice item
     * @param practiceItemId - The ID of the practice item being answered
     * @param userAnswer - The user's submitted answer
     * @returns Object containing evaluation results
     */
    async evaluateAnswer(practiceItemId: string, userAnswer: string): Promise<{ isCorrect: boolean }> {
        try {
            const response = await pb.send('/api/glimmer/v1/practice/evaluate-answer', {
                method: 'POST',
                body: {
                    practiceItemId,
                    userAnswer
                }
            });

            return response;
        } catch (err) {
            console.error('Failed to evaluate answer:', err);
            if (err instanceof Error) {
                throw err;
            }
            throw new Error('Failed to evaluate answer');
        }
    }

    /**
     * Processes a user's answer for a practice item (evaluation + scoring + database update)
     * @param practiceItemId - The ID of the practice item being answered
     * @param userAnswer - The user's submitted answer
     * @param practiceSession - The ID of the practice session
     * @param learnerId - The ID of the learner
     * @param hintLevelReached - The level of hints used (0 = no hints)
     * @returns Object containing complete processing results
     */
    async processAnswer(
        practiceItemId: string, 
        userAnswer: string, 
        practiceSession: string, 
        learnerId: string, 
        hintLevelReached: number = 0
    ): Promise<{ 
        isCorrect: boolean; 
        score: number; 
        feedback: string; 
        hintLevelReached: number; 
        attemptNumber: number 
    }> {
        try {
            const response = await pb.send('/api/glimmer/v1/practice/process-answer', {
                method: 'POST',
                body: {
                    practiceItemId,
                    userAnswer,
                    practiceSession,
                    learnerId,
                    hintLevelReached
                }
            });

            return response;
        } catch (err) {
            console.error('Failed to process answer:', err);
            if (err instanceof Error) {
                throw err;
            }
            throw new Error('Failed to process answer');
        }
    }
}

export const answersService = new AnswersService(); 