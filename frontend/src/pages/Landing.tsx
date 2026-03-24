import { Link } from 'react-router-dom';
import createProposalImg from '../assets/screenshot-create-proposal.png';
import weeklyPlanImg from '../assets/screenshot-weekly-plan.png';
import shoppingListImg from '../assets/screenshot-shoppinglist.png';

export function Landing() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-purple-600 via-pink-500 to-orange-400">
      {/* Header */}
      <header className="bg-white/10 backdrop-blur-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6 flex justify-between items-center">
          <div className="flex items-center space-x-2">
            <span className="text-4xl">🎲</span>
            <h1 className="text-3xl font-bold text-white">DishDice</h1>
          </div>
          <div className="space-x-4">
            <Link
              to="/login"
              className="text-white hover:text-gray-100 font-medium"
            >
              Log In
            </Link>
            <Link
              to="/register"
              className="bg-white text-purple-600 px-6 py-2 rounded-lg font-semibold hover:bg-gray-100 transition-colors"
            >
              Get Started
            </Link>
          </div>
        </div>
      </header>

      {/* Hero Section */}
      <section className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20 text-center">
        <h2 className="text-5xl md:text-6xl font-bold text-white mb-6">
          Stop wondering<br />
          <span className="text-yellow-300">"What should we eat this week?"</span>
        </h2>
        <p className="text-xl md:text-2xl text-white/90 mb-12 max-w-3xl mx-auto">
          Let AI create personalized 7-day meal plans with recipes and shopping lists.
          Save time, eat better, and never repeat the same meal.
        </p>
        <Link
          to="/register"
          className="inline-block bg-white text-purple-600 px-8 py-4 rounded-xl text-lg font-bold hover:bg-gray-100 transition-colors shadow-xl"
        >
          Start Planning Free
        </Link>
      </section>

      {/* How It Works */}
      <section className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        <h3 className="text-3xl font-bold text-white text-center mb-16">
          How DishDice Works
        </h3>

        <div className="space-y-24">
          {/* Step 1 */}
          <div className="flex flex-col md:flex-row items-center gap-12">
            <div className="md:w-1/2">
              <div className="bg-white rounded-xl shadow-2xl overflow-hidden">
                <img
                  src={createProposalImg}
                  alt="Create meal plan form"
                  className="w-full"
                />
              </div>
            </div>
            <div className="md:w-1/2 text-white">
              <div className="inline-block bg-white/20 rounded-full px-4 py-2 mb-4">
                <span className="font-bold">Step 1</span>
              </div>
              <h4 className="text-3xl font-bold mb-4">Set Your Preferences</h4>
              <p className="text-lg text-white/90 leading-relaxed">
                Tell us what you want to eat this week, what ingredients you already have,
                and any dietary restrictions. AI will consider all your preferences when
                generating your personalized meal plan.
              </p>
            </div>
          </div>

          {/* Step 2 */}
          <div className="flex flex-col md:flex-row-reverse items-center gap-12">
            <div className="md:w-1/2">
              <div className="bg-white rounded-xl shadow-2xl overflow-hidden">
                <img
                  src={weeklyPlanImg}
                  alt="Weekly meal plan with recipes"
                  className="w-full"
                />
              </div>
            </div>
            <div className="md:w-1/2 text-white">
              <div className="inline-block bg-white/20 rounded-full px-4 py-2 mb-4">
                <span className="font-bold">Step 2</span>
              </div>
              <h4 className="text-3xl font-bold mb-4">Get Your 7-Day Plan</h4>
              <p className="text-lg text-white/90 leading-relaxed">
                AI instantly generates a complete weekly meal plan with detailed recipes
                and ingredient lists. Don't like a day? Regenerate it for 3 new options.
                Add meals to your shopping list with one click.
              </p>
            </div>
          </div>

          {/* Step 3 */}
          <div className="flex flex-col md:flex-row items-center gap-12">
            <div className="md:w-1/2">
              <div className="bg-white rounded-xl shadow-2xl overflow-hidden">
                <img
                  src={shoppingListImg}
                  alt="Shopping list with checkboxes"
                  className="w-full"
                />
              </div>
            </div>
            <div className="md:w-1/2 text-white">
              <div className="inline-block bg-white/20 rounded-full px-4 py-2 mb-4">
                <span className="font-bold">Step 3</span>
              </div>
              <h4 className="text-3xl font-bold mb-4">Shop & Cook</h4>
              <p className="text-lg text-white/90 leading-relaxed">
                Your shopping list is automatically organized with all the ingredients you need.
                Check items off as you shop, then follow the recipes to cook delicious meals
                all week long.
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* Features */}
      <section className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20">
        <h3 className="text-3xl font-bold text-white text-center mb-16">
          Why Choose DishDice?
        </h3>
        <div className="grid md:grid-cols-3 gap-8">
          <div className="bg-white/10 backdrop-blur-sm rounded-xl p-8">
            <div className="text-4xl mb-4">🤖</div>
            <h4 className="text-xl font-bold text-white mb-3">AI-Powered</h4>
            <p className="text-white/90">
              Powered by OpenAI GPT-4, our AI creates personalized meal plans that match
              your taste and dietary needs.
            </p>
          </div>

          <div className="bg-white/10 backdrop-blur-sm rounded-xl p-8">
            <div className="text-4xl mb-4">🔄</div>
            <h4 className="text-xl font-bold text-white mb-3">No Repeats</h4>
            <p className="text-white/90">
              AI tracks your meal history to ensure variety. You'll never get the same
              meal twice in a row.
            </p>
          </div>

          <div className="bg-white/10 backdrop-blur-sm rounded-xl p-8">
            <div className="text-4xl mb-4">⚡</div>
            <h4 className="text-xl font-bold text-white mb-3">Save Time</h4>
            <p className="text-white/90">
              Stop spending hours planning meals and making shopping lists. Get your
              complete weekly plan in seconds.
            </p>
          </div>

          <div className="bg-white/10 backdrop-blur-sm rounded-xl p-8">
            <div className="text-4xl mb-4">🎯</div>
            <h4 className="text-xl font-bold text-white mb-3">Personalized</h4>
            <p className="text-white/90">
              Set your dietary preferences once and get meal plans tailored to your
              specific needs every week.
            </p>
          </div>

          <div className="bg-white/10 backdrop-blur-sm rounded-xl p-8">
            <div className="text-4xl mb-4">🛒</div>
            <h4 className="text-xl font-bold text-white mb-3">Smart Shopping</h4>
            <p className="text-white/90">
              Automatically generated shopping lists with precise quantities. Add meals
              to your list with one click.
            </p>
          </div>

          <div className="bg-white/10 backdrop-blur-sm rounded-xl p-8">
            <div className="text-4xl mb-4">💰</div>
            <h4 className="text-xl font-bold text-white mb-3">Free to Start</h4>
            <p className="text-white/90">
              Get started with DishDice today at no cost. Plan better meals for your
              family without breaking the bank.
            </p>
          </div>
        </div>
      </section>

      {/* CTA */}
      <section className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-20 text-center">
        <h3 className="text-4xl font-bold text-white mb-6">
          Ready to Transform Your Meal Planning?
        </h3>
        <p className="text-xl text-white/90 mb-8">
          Join DishDice today and never stress about "what's for dinner" again.
        </p>
        <Link
          to="/register"
          className="inline-block bg-white text-purple-600 px-10 py-4 rounded-xl text-lg font-bold hover:bg-gray-100 transition-colors shadow-xl"
        >
          Get Started Now
        </Link>
      </section>

      {/* Footer */}
      <footer className="bg-black/20 backdrop-blur-sm mt-20">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 text-center text-white/70">
          <p>&copy; 2026 DishDice. AI-powered meal planning made simple.</p>
        </div>
      </footer>
    </div>
  );
}
