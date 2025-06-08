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
}

export const answersService = new AnswersService(); 